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
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil{
		return err
	}
	fmt.Println(feed)
	return nil
}

func addFeedHandler(s *State, args []string) error{
	if len(args) < 2{
		return fmt.Errorf("Missing %v argument(s)\nUsage: addfeed <name> <url>", 2 - len(args))
	}
	user, err := s.DB.GetUserByName(context.Background(), s.Config.UserName)
	if err != nil{
		return err
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
	err = followHandler(s, []string{feed.FeedUrl})
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

func followHandler(s *State, args []string) error{
	if len(args) < 1{
		return fmt.Errorf("Missing argument feed url\nUsage: follow <feed_url>")
	}
	feed, err := s.DB.GetFeedByUrl(context.Background(), args[0])
	if err != nil{
		return fmt.Errorf("Error finding feed: %w", err)
	} 
	user, err := s.DB.GetUserByName(context.Background(), s.Config.UserName)
	if err != nil{
		return fmt.Errorf("Error looking up current user: %w", err)
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

func followingHandler(s *State, args []string) error{
	user, err := s.DB.GetUserByName(context.Background(), s.Config.UserName)
	if err != nil{
		return fmt.Errorf("Error looking up current user: %w", err)
	}
	feeds, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil{
		return fmt.Errorf("Error getting feed follows for user: %w", err)
	}
	for _, feed := range feeds{
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}