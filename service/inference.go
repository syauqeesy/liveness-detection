package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/liveness-detection/payload"
	"github.com/liveness-detection/protobuf/compiled/inference"
)

type InferenceService interface {
	Execute(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*payload.ExecuteInferenceResponse, error)
}

type inferenceService service

func (s *inferenceService) Execute(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*payload.ExecuteInferenceResponse, error) {
	imageInBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	response, err := s.GRPCOutbound.Inference.Predict(ctx, &inference.PredictionRequest{
		Image: imageInBytes,
	})
	if err != nil {
		return nil, err
	}

	return &payload.ExecuteInferenceResponse{
			Result: response.GetResult(),
			Live:   response.GetLive(),
			Spoof:  response.GetSpoof(),
		},
		nil
}
