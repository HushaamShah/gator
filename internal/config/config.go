package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBURL string `json:"db_url"`
	User  string `json:"user"`
}

var config Config

func Read() (Config, error) {
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return Config{}, err
	}

	content, err := os.ReadFile(UserHomeDir + "/.gatorconfig.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return Config{}, err
	}
	err1 := json.Unmarshal(content, &config)
	if err1 != nil {
		fmt.Println("Error unmarshalling config file:", err1)
		return Config{}, err1
	}
	return config, nil
}

func (c *Config) SetUser(name string) {
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	config.User = name
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling config file:", err)
		return
	}

	err = os.WriteFile(UserHomeDir+"/.gatorconfig.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}
	fmt.Println("Config file updated successfully")
}
