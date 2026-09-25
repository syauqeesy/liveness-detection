package service

import (
	"github.com/syauqeesy/liveness-detection/configuration"
	grpc_outbound "github.com/syauqeesy/liveness-detection/outbound/grpc"
)

type service struct {
	Configuration *configuration.Configuration
	GRPCOutbound  *grpc_outbound.GRPCOutboundService
}

type Service struct {
	Inference InferenceService
}

func NewService(configuration *configuration.Configuration, grpcOutbound *grpc_outbound.GRPCOutboundService) *Service {
	svc := &service{
		Configuration: configuration,
		GRPCOutbound:  grpcOutbound,
	}

	return &Service{
		Inference: (*inferenceService)(svc),
	}
}
