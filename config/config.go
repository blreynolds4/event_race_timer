package config

import (
	"encoding/json"
)

type RaceConfig struct {
	RaceName     string
	RedisAddress string
	RedisDbumber int
	SourceRanks  map[string]int
}

func LoadConfigData(data []byte) (RaceConfig, error) {
	var raceConfig RaceConfig
	err := json.Unmarshal([]byte(data), &raceConfig)
	if err != nil {
		return raceConfig, err
	}

	return raceConfig, err
}

func GetConfigData(rc RaceConfig) ([]byte, error) {
	return json.MarshalIndent(rc, "", "  ")
}
