package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Application struct {
		Service     string `json:"service"`
		Environment string `json:"environment"`
		Client      string `json:"client"`
	}
	Http struct {
		Port string `json:"port"`
	} `json:"http"`
	GRPC struct {
		Service struct {
			Detector string `json:"detector"`
		} `json:"service"`
	} `json:"grpc"`
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
