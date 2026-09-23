package grpc_outbound

import (
	"github.com/liveness-detection/configuration"
	"github.com/liveness-detection/protobuf/compiled/inference"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCOutboundConnection struct {
	DetectorService *grpc.ClientConn
}

type GRPCOutboundService struct {
	Inference inference.InferenceServiceClient
}

func New(configuration *configuration.Configuration) *GRPCOutboundService {
	detectorServiceConnection, err := grpc.NewClient(configuration.GRPC.Service.Detector, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	grpcOutboundService := &GRPCOutboundService{
		Inference: inference.NewInferenceServiceClient(detectorServiceConnection),
	}

	return grpcOutboundService
}

func (o *GRPCOutboundConnection) Close() error {
	if o.DetectorService == nil {
		return nil
	}

	o.DetectorService.Close()

	return nil
}
