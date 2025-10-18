package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type RaceConfig struct {
	RaceName      string
	RedisAddress  string
	RedisDbNumber int
	PgConnect     string
	SourceRanks   map[string]int
}

func LoadConfigData(configPath string, raceConfig *RaceConfig) error {
	fmt.Println("Loading config from", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileContent, raceConfig)
}

func GetConfigData(rc RaceConfig) ([]byte, error) {
	return json.MarshalIndent(rc, "", "  ")
}
