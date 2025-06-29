package command

import (
	"context"
	"fmt"
	"math/rand"
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

func addFeedHandler(s *State, args []string, user database.User) error{
	if len(args) < 2{
		return fmt.Errorf("Missing %v argument(s)\nUsage: addfeed <name> <url>", 2 - len(args))
	}
	params := database.CreateFeedParams{
		ID: rand.Int31(),
		CreatedAt: time.Now(),
		FeedName: args[0],
		FeedUrl: args[1],
		UserID: user.ID,
	}
	feed, err := s.DB.CreateFeed(context.Background(), params)
	if err != nil{
		return fmt.Errorf("Error creating feed: %w", err)
	}
	err = followHandler(s, []string{feed.FeedUrl}, user)
	return err
}

func feedsHandler(s *State, args []string) error{
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil{
		return fmt.Errorf("Error retrieving feed data: %w", err)
	}
	for _, feed := range feeds{
		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)
		if err != nil{
			return fmt.Errorf("Error looking up user with id %v: %w", feed.UserID, err)
		}
		fmt.Printf("* %s:%s added by %s\n", feed.FeedName, feed.FeedUrl, user.UserName)
	}
	return nil
}

func followHandler(s *State, args []string, user database.User) error{
	if len(args) < 1{
		return fmt.Errorf("Missing argument feed url\nUsage: follow <feed_url>")
	}
	feed, err := s.DB.GetFeedByUrl(context.Background(), args[0])
	if err != nil{
		return fmt.Errorf("Error finding feed: %w", err)
	}
	params := database.CreateFeedFollowParams{
		ID: rand.Int31(),
		CreatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), params)
	if err != nil{
		return fmt.Errorf("Error creating feed follow row: %w", err)
	}
	fmt.Printf("%s follows %s", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func followingHandler(s *State, args []string, user database.User) error{
	feeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil{
		return fmt.Errorf("Error getting feed follows for user: %w", err)
	}
	for _, feed := range feeds{
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}

func unFollowHandler(s *State, args []string, user database.User) error{
	if len(args) < 1{
		return fmt.Errorf("Missing feed url argument\nUsage: unfollow <feed_url>")
	}
	err := s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, FeedUrl: args[0]})
	if err != nil{
		return fmt.Errorf("Error unfollowing feed %s: %w", args[0], err)
	}
	fmt.Printf("%s unfollows %s", user.UserName, args[0])
	return nil
}