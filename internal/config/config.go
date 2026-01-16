package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".pdbmdt.json"

type Config struct {
	Db_URL string `json:"db_url"`
}

func Read() (Config, error) {
	home_directory, err := GetHomeDirectory()
	if err != nil {
		return Config{}, err
	}
	file_path := home_directory + "/" + configFileName
	var config Config
	file, err := os.Open(file_path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("Not decoding\n")
		return Config{}, err
	}
	return config, nil
}
func GetHomeDirectory() (string, error) {
	return os.UserHomeDir()
}
