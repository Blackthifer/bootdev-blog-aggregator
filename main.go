package main

import (
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"github.com/Blackthifer/bootdev-blog-aggregator/internal/command"
	"github.com/Blackthifer/bootdev-blog-aggregator/internal/config"
)

const dbConnStr = "postgres://postgres:postgres@localhost:5432/gator"

func main(){
	gatorConfig, err := config.Read()
	if err != nil{
		fmt.Println(err)
	}
	state := &command.State{
		Config: gatorConfig,
	}
	cmds := command.InitCommands()
	cmd := os.Args[1:]
	if len(cmd) == 0{
		fmt.Println("No command specified")
		os.Exit(1)
	}
	err = cmds.Run(state, cmd[0], cmd[1:])
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}