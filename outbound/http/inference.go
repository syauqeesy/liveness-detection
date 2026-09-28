package outbound_http

import (
	"context"
	"fmt"

	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
)

type ManagedInference interface {
	Predict(ctx context.Context, image string) (Prediction, error)
}

type managedInference struct {
	configuration *configuration.Configuration
	httpClient    common.CommonHttpClient
	logger        common.Logger
}

func NewManagedInference(
	configuration *configuration.Configuration,
	httpClient common.CommonHttpClient,
	logger common.Logger,
) ManagedInference {
	return &managedInference{
		configuration: configuration,
		httpClient:    httpClient,
		logger:        logger,
	}
}

type managedPredictionRequest struct {
	Instances []managedPredictionInstance `json:"instances"`
}

type managedPredictionInstance struct {
	RequestId string `json:"request_id"`
	Image     string `json:"image"`
}

type managedPredictionResponse struct {
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Result string  `json:"result"`
	Live   float32 `json:"live"`
	Spoof  float32 `json:"spoof"`
}

func (m *managedInference) Predict(
	ctx context.Context,
	image string,
) (Prediction, error) {
	requestId := common.RequestIdFromContext(ctx)
	request := managedPredictionRequest{
		Instances: []managedPredictionInstance{
			{
				RequestId: requestId,
				Image:     image,
			},
		},
	}

	response := managedPredictionResponse{}

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
