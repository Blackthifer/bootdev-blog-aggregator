package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Blackthifer/bootdev-blog-aggregator/internal/command"
	"github.com/Blackthifer/bootdev-blog-aggregator/internal/config"
	"github.com/Blackthifer/bootdev-blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

const dbConnStr = "postgres://postgres:postgres@localhost:5432/gator"

func main(){
	gatorConfig, err := config.Read()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	dbQ := database.New(db)
	state := &command.State{
		Config: gatorConfig,
		DB: dbQ,
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