package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Application struct {
		Service     string `json:"service"`
		Environment string `json:"environment"`
	}
	Http struct {
		Port string `json:"port"`
	} `json:"http"`
}

func NewConfiguration(path string) (*Configuration, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var configuration Configuration

	err = json.Unmarshal(bytes, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}
