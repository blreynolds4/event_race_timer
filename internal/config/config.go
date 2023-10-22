package config

import (
	"encoding/json"
	"io/ioutil"
)

type RaceConfig struct {
	// RaceName     string
	// RedisAddress string
	// RedisDbumber int
	SourceRanks map[string]int
}

func LoadConfigData(configPath string, raceConfig *RaceConfig) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(file), raceConfig)
}

func GetConfigData(rc RaceConfig) ([]byte, error) {
	return json.MarshalIndent(rc, "", "  ")
}
