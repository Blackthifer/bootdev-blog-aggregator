package command

import (
	"fmt"

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
	err := s.Config.SetUser(args[0])
	if err != nil{
		return fmt.Errorf("login failed: %w", err)
	}
	fmt.Printf("User was set to %s\n", args[0])
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
	return cmds
}