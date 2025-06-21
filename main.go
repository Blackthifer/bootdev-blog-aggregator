package main

import (
	"fmt"

	"github.com/Blackthifer/bootdev-blog-aggregator/internal/config"
)

func main(){
	gatorConfig := config.Read()
	gatorConfig.SetUser("Hessel")
	gatorConfig = config.Read()
	fmt.Println("Config:")
	fmt.Println("- db_url:", gatorConfig.DbUrl)
	fmt.Println("- current_user_name:", gatorConfig.UserName)
}