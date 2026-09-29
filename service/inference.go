package service

import (
	"context"
	"encoding/base64"
	"io"
	"mime/multipart"
	"time"

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

	requestId := common.RequestIdFromContext(ctx)

	switch mode {
	case "self_managed_service":
		started := time.Now()

		result, err := s.GRPCOutbound.Inference.Predict(ctx, &inference.PredictionRequest{
			RequestId: requestId,
			Image:     imageInBase64,
		})

		duration := time.Since(started)
		s.Logger.Info(
			"grpc request completed",
			"request_id",
			requestId,
			"duration_ms",
			duration.Seconds()*1000,
		)

		if err != nil {
			return nil, err
		}

		response.Result = result.GetResult()
		response.Live = result.GetLive()
		response.Spoof = result.GetSpoof()
	case "managed_service":
		inference := outbound_http.NewInference(s.Configuration, common.NewHttpClient(s.Logger), s.Logger)
		result, err := inference.Predict(ctx, imageInBase64)
		if err != nil {
			return nil, err
		}

		response.Result = result.Result
		response.Live = result.Live
		response.Spoof = result.Spoof
	}

	return response, nil
}
