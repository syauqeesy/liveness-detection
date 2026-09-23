package foundation

import (
	"context"
	"errors"

	"github.com/liveness-detection/common"
	"github.com/liveness-detection/configuration"
)

type Foundation interface {
	Setup() error
	Boot() error
	Shutdown(ctx context.Context) error
}

const (
	FoundationHttp = "http"
)

func NewFoundation(foundationType string, arguments []string, logger common.Logger, config *configuration.Configuration) (Foundation, error) {
	switch foundationType {
	case FoundationHttp:
		return &httpFoundation{
			configuration: config,
			logger:        logger,
		}, nil

	default:
		return nil, errors.New("invalid foundation type")
	}
}
