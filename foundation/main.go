package foundation

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
)

const (
	shutdownTimeout = 10 * time.Second

	FoundationHttp = "http"
)

type Foundation interface {
	Setup() error
	Boot() error
	Shutdown(ctx context.Context) error
}

type mainFoundation struct {
	foundation Foundation
	logger     common.Logger
}

func NewFoundation(
	foundationType string,
	arguments []string,
	logger common.Logger,
	configuration *configuration.Configuration,
) (Foundation, error) {
	switch foundationType {
	case FoundationHttp:
		return &mainFoundation{
			foundation: &httpFoundation{
				configuration: configuration,
				logger:        logger,
			},
			logger: logger,
		}, nil

	default:
		return nil, errors.New("invalid foundation type")
	}
}

func (f *mainFoundation) Setup() error {
	return f.foundation.Setup()
}

func (f *mainFoundation) Boot() error {
	bootErr := make(chan error, 1)

	go func() {
		bootErr <- f.foundation.Boot()
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-bootErr:
		return err

	case <-ctx.Done():
		f.logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := f.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func (f *mainFoundation) Shutdown(ctx context.Context) error {
	return f.foundation.Shutdown(ctx)
}

func Boot(arguments []string) error {
	if len(arguments) < 1 {
		return errors.New("foundation type is required")
	}

	configPath, err := filepath.Abs("./config.json")
	if err != nil {
		return err
	}

	config, err := configuration.NewConfiguration(configPath)
	if err != nil {
		return err
	}

	logger := common.NewLogger(
		config.Application.Service,
		config.Application.Environment,
	)

	foundation, err := NewFoundation(
		arguments[0],
		arguments[1:],
		logger,
		config,
	)
	if err != nil {
		logger.Error(
			"failed to create foundation",
			"error", err,
		)

		return err
	}

	if err := foundation.Setup(); err != nil {
		logger.Error(
			"failed to setup foundation",
			"error", err,
		)

		return err
	}

	if err := foundation.Boot(); err != nil {
		logger.Error(
			"foundation boot failed",
			"error", err,
		)

		return err
	}

	logger.Info("application shutdown completed")

	return nil
}
