package service

import "ahmadsyauqi.dev/projects/liveness-detection/configuration"

type service struct {
	Configuration *configuration.Configuration
}

type Service struct {
	Inference InferenceService
}

func NewService(configuration *configuration.Configuration) *Service {
	svc := &service{
		Configuration: configuration,
	}

	return &Service{
		Inference: (*inferenceService)(svc),
	}
}
