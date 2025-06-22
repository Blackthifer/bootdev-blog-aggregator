package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct{
	DbUrl string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const gatorConfig = ".gatorconfig.json"

func Read() (*Config, error){
	gatorPath, err := getConfigPath()
	if err != nil{
		return nil, fmt.Errorf("%w", err)
	}
	jsonData, err := os.ReadFile(gatorPath)
	if err != nil{
		return nil, fmt.Errorf("Error reading gatorconfig: %w", err)
	}
	jsonConfig := &Config{}
	if err = json.Unmarshal(jsonData, jsonConfig); err != nil{
		return nil, fmt.Errorf("Error parsing from json: %w", err)
	}
	return jsonConfig, nil
}

func (c *Config) SetUser(user string) error{
	c.UserName = user
	jsonData, err := json.Marshal(c)
	if err != nil{
		return fmt.Errorf("Error parsing into json: %w", err)
	}
	gatorPath, err := getConfigPath()
	if err != nil{
		return fmt.Errorf("%w", err)
	}
	err = os.WriteFile(gatorPath, jsonData, os.ModePerm)
	if err != nil{
		return fmt.Errorf("Error writing to file: %w", err)
	}
	return nil
}

func getConfigPath() (string, error){
	homeDir, err := os.UserHomeDir()
	if err != nil{
		return "", fmt.Errorf("Error looking up home dir: %w", err)
	}
	return homeDir + "/" + gatorConfig, nil
}