package command

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"math/rand"
	"strconv"
	"strings"
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
			fmt.Println(err)
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
	fmt.Printf("Fetched feed '%s':\n", dbFeed.FeedName)
	for _, item := range feed.Channel.Item{
		err = storePost(s, item, dbFeed.ID)
		if err != nil{
			fmt.Println(err)
		}
	}
	return nil
}

func storePost(s *State, item rss.RSSItem, feed_id int32) error{
	maybeDescription := sql.NullString{
		String: html.UnescapeString(item.Description),
		Valid: true,
	}
	if item.Description == ""{
		maybeDescription.Valid = false
	}
	pubTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
	if err != nil{
		return fmt.Errorf("Error parsing time: %w", err)
	}
	params := database.CreatePostParams{
		ID: rand.Int31(),
		CreatedAt: time.Now(),
		Title: item.Title,
		PostUrl: item.Link,
		PostDescription: maybeDescription,
		PublishedAt: pubTime,
		FeedID: feed_id,
	}
	_, err = s.DB.CreatePost(context.Background(), params)
	if err != nil{
		//fmt.Println(err)
		if !strings.Contains(err.Error(), "post_url"){
			return fmt.Errorf("Error creating post in database: %w", err)
		}
	}
	//fmt.Printf("* %s : %s\n", item.Title, item.PubDate)
	return nil
}

func browseHandler(s *State, args []string, user database.User) error{
	limit := 2
	if len(args) >= 1{
		argLimit, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil{
			return fmt.Errorf("Limit argument is not an integer")
		}
		limit = int(argLimit)
	}
	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil{
		return fmt.Errorf("Error looking up posts: %w", err)
	}
	for _, post := range posts{
		fmt.Println("-------------------------------------")
		fmt.Printf("%s :: %s\n", post.Title, post.PostUrl)
		if post.PostDescription.Valid{
			fmt.Println(post.PostDescription.String)
		}
	}
	return nil
}