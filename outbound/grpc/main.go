package grpc_outbound

import (
	"github.com/syauqeesy/liveness-detection/configuration"
	"github.com/syauqeesy/liveness-detection/pb/compiled/inference"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCOutboundConnection struct {
	PredictionService *grpc.ClientConn
}

type GRPCOutboundService struct {
	Inference inference.InferenceServiceClient
}

func New(configuration *configuration.Configuration) *GRPCOutboundService {
	detectorServiceConnection, err := grpc.NewClient(configuration.Service.SelfManagedServicePrediction.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	grpcOutboundService := &GRPCOutboundService{
		Inference: inference.NewInferenceServiceClient(detectorServiceConnection),
	}

	return grpcOutboundService
}

func (o *GRPCOutboundConnection) Close() error {
	if o.PredictionService == nil {
		return nil
	}

	o.PredictionService.Close()

	return nil
}
