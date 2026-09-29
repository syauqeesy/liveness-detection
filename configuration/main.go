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
	} `json:"application"`

	Http struct {
		Port string `json:"port"`
	} `json:"http"`

	Service struct {
		SelfManagedServicePrediction struct {
			Endpoint string `json:"endpoint"`
		} `json:"self_managed_service_prediction"`

		ManagedServicePrediction struct {
			Endpoint           string `json:"endpoint"`
			AuthorizationToken string `json:"authorization_token"`
		} `json:"managed_service_prediction"`
	} `json:"service"`
}

func NewConfiguration(path string) (*Configuration, error) {
	var configuration Configuration

	if bytes, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(bytes, &configuration); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if value := os.Getenv("APPLICATION_SERVICE"); value != "" {
		configuration.Application.Service = value
	}

	if value := os.Getenv("APPLICATION_ENVIRONMENT"); value != "" {
		configuration.Application.Environment = value
	}

	if value := os.Getenv("APPLICATION_CLIENT"); value != "" {
		configuration.Application.Client = value
	}

	if value := os.Getenv("HTTP_PORT"); value != "" {
		configuration.Http.Port = value
	}

	if value := os.Getenv("SERVICE_SELF_MANAGED_SERVICE_PREDICTION_ENDPOINT"); value != "" {
		configuration.Service.SelfManagedServicePrediction.Endpoint = value
	}

	if value := os.Getenv("SERVICE_MANAGED_SERVICE_PREDICTION_ENDPOINT"); value != "" {
		configuration.Service.ManagedServicePrediction.Endpoint = value
	}

	if value := os.Getenv("SERVICE_MANAGED_SERVICE_PREDICTION_AUTHORIZATION_TOKEN"); value != "" {
		configuration.Service.ManagedServicePrediction.AuthorizationToken = value
	}

	return &configuration, nil
}
