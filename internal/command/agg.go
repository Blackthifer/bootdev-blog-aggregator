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
		return fmt.Errorf("Missing %v arguments", 2 - len(args))
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
	_, err = s.DB.CreateFeed(context.Background(), params)
	if err != nil{
		return fmt.Errorf("Error creating feed: %w", err)
	}
	return nil
}