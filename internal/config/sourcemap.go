package config

import (
	"encoding/json"
	"io/ioutil"
)

type SourceConfig struct {
	// map is hostname to source name
	SourceMap map[string]string `json:"sources"`
}

func LoadAnyConfigData[CT any](configPath string, c *CT) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(file), c)
}
