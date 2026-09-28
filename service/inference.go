package service

import (
	"context"
	"encoding/base64"
	"io"
	"mime/multipart"

	"github.com/syauqeesy/liveness-detection/common"
	outbound_http "github.com/syauqeesy/liveness-detection/outbound/http"
	"github.com/syauqeesy/liveness-detection/payload"
	"github.com/syauqeesy/liveness-detection/pb/compiled/inference"
)

type InferenceService interface {
	Predict(ctx context.Context, mode string, file multipart.File, header *multipart.FileHeader) (*payload.ExecuteInferenceResponse, error)
}

type inferenceService service

func (s *inferenceService) Predict(ctx context.Context, mode string, file multipart.File, header *multipart.FileHeader) (*payload.ExecuteInferenceResponse, error) {
	response := &payload.ExecuteInferenceResponse{}

	imageInBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	imageInBase64 := base64.StdEncoding.EncodeToString(imageInBytes)

	switch mode {
	case "self_managed_service":
		result, err := s.GRPCOutbound.Inference.Predict(ctx, &inference.PredictionRequest{
			RequestId: common.RequestIdFromContext(ctx),
			Image:     imageInBase64,
		})
		if err != nil {
			return nil, err
		}

		response.Result = result.GetResult()
		response.Live = result.GetLive()
		response.Spoof = result.GetSpoof()
	case "managed_service":
		managedInference := outbound_http.NewManagedInference(s.Configuration, common.NewHttpClient(s.Logger), s.Logger)
		result, err := managedInference.Predict(ctx, imageInBase64)
		if err != nil {
			return nil, err
		}

		response.Result = result.Result
		response.Live = result.Live
		response.Spoof = result.Spoof
	}

	return response, nil
}
