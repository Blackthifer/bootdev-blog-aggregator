package command

import (
	"context"
	"fmt"

	"github.com/Blackthifer/bootdev-blog-aggregator/internal/database"
)

func requireLoggedInUser(handler func(s *State, args []string, user database.User) error) func(*State, []string) error{
	loggedInHandler := func(s *State, args []string) error{
		user, err := s.DB.GetUserByName(context.Background(), s.Config.UserName)
		if err != nil{
			return fmt.Errorf("Error looking up current user: %w", err)
		}
		return handler(s, args, user)
	}
	return loggedInHandler
}