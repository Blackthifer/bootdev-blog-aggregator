package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct{
	DbUrl string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const gatorConfig = ".gatorconfig.json"

func Read() *Config{
	gatorPath, err := getConfigPath()
	if err != nil{
		log.Println(err)
		return nil
	}
	jsonData, err := os.ReadFile(gatorPath)
	if err != nil{
		log.Println("Error reading gatorconfig:", err)
		return nil
	}
	jsonConfig := &Config{}
	if err = json.Unmarshal(jsonData, jsonConfig); err != nil{
		log.Println("Error parsing from json:", err)
		return nil
	}
	return jsonConfig
}

func (c *Config) SetUser(user string){
	c.UserName = user
	jsonData, err := json.Marshal(c)
	if err != nil{
		log.Println("Error parsing into json:", err)
		return
	}
	gatorPath, err := getConfigPath()
	if err != nil{
		log.Println(err)
		return
	}
	err = os.WriteFile(gatorPath, jsonData, os.ModePerm)
	if err != nil{
		log.Println("Error writing to file:", err)
	}
}

func getConfigPath() (string, error){
	homeDir, err := os.UserHomeDir()
	if err != nil{
		return "", fmt.Errorf("Error looking up home dir: %w", err)
	}
	return homeDir + "/" + gatorConfig, nil
}