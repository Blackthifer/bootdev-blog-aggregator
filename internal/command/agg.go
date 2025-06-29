package command

import (
	"context"
	"fmt"
	"time"

	"github.com/Blackthifer/bootdev-blog-aggregator/internal/database"
	"github.com/Blackthifer/bootdev-blog-aggregator/internal/rss"
)

func aggHandler(s *State, args []string) error{
	if len(args) < 1{
		return fmt.Errorf("Missing scrape delay argument\nUsage: agg <scrape_delay>")
	}
	delay, err := time.ParseDuration(args[0])
	if err != nil{
		return fmt.Errorf("Invalid delay format: %w", err)
	}
	for ticker := time.NewTicker(delay); ; <-ticker.C{
		err = scrapeFeed(s)
		if err != nil{
			return err
		}
	}
	//return nil
}

func scrapeFeed(s *State) error{
	dbFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil{
		return fmt.Errorf("Error fetching oldest feed data: %w", err)
	}
	mark, err := s.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: dbFeed.ID, UpdatedAt: time.Now()})
	if err != nil{
		return fmt.Errorf("Error marking feed as fetched: %w", err)
	}
	if !mark.LastFetchedAt.Valid || mark.UpdatedAt.Compare(mark.LastFetchedAt.Time) != 0{
		return fmt.Errorf("Feed at %s was not updated correctly", mark.FeedUrl)
	}
	feed, err := rss.FetchFeed(context.Background(), mark.FeedUrl)
	if err != nil{
		return err
	}
	fmt.Printf("Feed %s:\n", dbFeed.FeedName)
	for _, item := range feed.Channel.Item{
		fmt.Printf("* %s : %s\n", item.Title, item.PubDate)
	}
	return nil
}