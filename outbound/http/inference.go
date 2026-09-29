package outbound_http

import (
	"context"
	"fmt"

	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
)

type Inference interface {
	Predict(ctx context.Context, image string) (Prediction, error)
}

type inference struct {
	configuration *configuration.Configuration
	httpClient    common.CommonHttpClient
	logger        common.Logger
}

func NewInference(
	configuration *configuration.Configuration,
	httpClient common.CommonHttpClient,
	logger common.Logger,
) Inference {
	return &inference{
		configuration: configuration,
		httpClient:    httpClient,
		logger:        logger,
	}
}

type predictionRequest struct {
	Instances []predictionInstance `json:"instances"`
}

type predictionInstance struct {
	RequestId string `json:"request_id"`
	Image     string `json:"image"`
}

type predictionResponse struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Result string  `json:"result"`
	Live   float32 `json:"live"`
	Spoof  float32 `json:"spoof"`
}

func (m *inference) Predict(
	ctx context.Context,
	image string,
) (Prediction, error) {
	requestId := common.RequestIdFromContext(ctx)
	request := predictionRequest{
		Instances: []predictionInstance{
			{
				RequestId: requestId,
				Image:     image,
			},
		},
	}

	response := predictionResponse{}

	err := m.httpClient.PostJson(
		ctx,
		m.configuration.Service.ManagedServicePrediction.Endpoint,
		map[string]string{
			"Authorization": "Bearer " + m.configuration.Service.ManagedServicePrediction.AuthorizationToken,
			"X-Request-Id":  requestId,
		},
		request,
		&response,
	)
	if err != nil {
		m.logger.Error(
			"managed inference request failed",
			"request_id", requestId,
			"error", err,
		)

		return Prediction{}, fmt.Errorf(
			"managed inference request failed: %w",
			err,
		)
	}

	if len(response.Predictions) != 1 {
		err := fmt.Errorf(
			"expected 1 prediction, got %d",
			len(response.Predictions),
		)

		m.logger.Error(
			"managed inference returned invalid response",
			"request_id", requestId,
			"error", err,
		)

		return Prediction{}, err
	}

	return response.Predictions[0], nil
}
