package service

import (
	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
	grpc_outbound "github.com/syauqeesy/liveness-detection/outbound/grpc"
)

type service struct {
	Configuration *configuration.Configuration
	Logger        common.Logger
	GRPCOutbound  *grpc_outbound.GRPCOutboundService
}

type Service struct {
	Inference InferenceService
}

func NewService(configuration *configuration.Configuration, logger common.Logger, grpcOutbound *grpc_outbound.GRPCOutboundService) *Service {
	svc := &service{
		Configuration: configuration,
		Logger:        logger,
		GRPCOutbound:  grpcOutbound,
	}

	return &Service{
		Inference: (*inferenceService)(svc),
	}
}
