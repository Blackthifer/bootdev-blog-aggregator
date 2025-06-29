package command

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Blackthifer/bootdev-blog-aggregator/internal/config"
	"github.com/Blackthifer/bootdev-blog-aggregator/internal/database"
)

type State struct{
	Config *config.Config
	DB *database.Queries
}

type Command struct{
	Name string
	Args []string
}

type Commands struct{
	handlers map[string]func(*State, []string)error
}

func loginHandler(s *State, args []string) error{
	if len(args) == 0{
		return fmt.Errorf("login is missing username argument")
	}
	_, err := s.DB.GetUserByName(context.Background(), args[0])
	if err != nil{
		return fmt.Errorf("user '%s' does not exist", args[0])
	}
	err = s.Config.SetUser(args[0])
	if err != nil{
		return fmt.Errorf("login failed: %w", err)
	}
	fmt.Printf("User was set to %s\n", args[0])
	return nil
}

func registerHandler(s *State, args []string) error{
	if len(args) == 0{
		return fmt.Errorf("register is missing username argument")
	}
	_, err := s.DB.GetUserByName(context.Background(), args[0])
	if err == nil{
		return fmt.Errorf("user '%s' already exists", args[0])
	}
	uParams := database.CreateUserParams{
		ID: rand.Int31(),
		CreatedAt: time.Now(),
		UserName: args[0],
	}
	_, err = s.DB.CreateUser(context.Background(), uParams)
	if err != nil{
		return fmt.Errorf("Error creating user: %w", err)
	}
	log.Println(uParams)
	fmt.Println("Created user: ", args[0])
	err = loginHandler(s, args)
	if err != nil{
		return fmt.Errorf("%w", err)
	}
	return nil
}

func resetHandler(s *State, args []string) error{
	err := s.DB.DeleteAllUsers(context.Background())
	if err != nil{
		return fmt.Errorf("Reset failed: %w", err)
	}
	fmt.Println("Reset complete")
	return nil
}

func usersHandler(s *State, args []string) error{
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed fetching users: %w", err)
	}
	for _, user := range users{
		printStr := "* " + user.UserName
		if user.UserName == s.Config.UserName{
			printStr += " (current)"
		}
		fmt.Println(printStr)
	}
	return nil
}

func (c *Commands) Run(s *State, cmd string, args []string) error{
	cmdFunc, ok := c.handlers[cmd]
	if !ok{
		return fmt.Errorf("unknown command: %s", cmd)
	}
	err := cmdFunc(s, args)
	if err != nil{
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (c *Commands) register(name string, f func(*State, []string) error){
	c.handlers[name] = f
}

func InitCommands() *Commands{
	cmds := &Commands{
		handlers: map[string]func(*State, []string) error{},
	}
	cmds.register("login", loginHandler)
	cmds.register("register", registerHandler)
	cmds.register("reset", resetHandler)
	cmds.register("users", usersHandler)
	cmds.register("agg", aggHandler)
	cmds.register("addfeed", requireLoggedInUser(addFeedHandler))
	cmds.register("feeds", feedsHandler)
	cmds.register("follow", requireLoggedInUser(followHandler))
	cmds.register("following", requireLoggedInUser(followingHandler))
	cmds.register("unfollow", requireLoggedInUser(unFollowHandler))
	return cmds
}