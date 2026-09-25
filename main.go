package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
	"github.com/syauqeesy/liveness-detection/foundation"
)

const (
	shutdownTimeout = 10 * time.Second
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	config, err := loadConfiguration()
	if err != nil {
		os.Exit(1)
	}

	logger := common.NewLogger(
		config.Application.Service,
		config.Application.Environment,
	)

	f, err := createFoundation(
		config,
		logger,
		os.Args[1],
		os.Args[2:],
	)
	if err != nil {
		logger.Error(
			"failed to create foundation",
			"error", err,
		)

		os.Exit(1)
	}

	if err := setupFoundation(f, logger); err != nil {
		os.Exit(1)
	}

	if err := bootFoundation(f, logger); err != nil {
		os.Exit(1)
	}
}

func loadConfiguration() (*configuration.Configuration, error) {
	configPath, err := filepath.Abs("./config.json")
	if err != nil {
		return nil, err
	}

	return configuration.NewConfiguration(configPath)
}

func createFoundation(config *configuration.Configuration, logger common.Logger, foundationType string, arguments []string) (foundation.Foundation, error) {
	return foundation.NewFoundation(
		foundationType,
		arguments,
		logger,
		config,
	)
}

func setupFoundation(f foundation.Foundation, logger common.Logger) error {
	if err := f.Setup(); err != nil {
		logger.Error(
			"failed to setup foundation",
			"error", err,
		)

		return err
	}

	return nil
}

func bootFoundation(f foundation.Foundation, logger common.Logger) error {
	bootErr := make(chan error, 1)

	go func() {
		bootErr <- f.Boot()
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-bootErr:
		if err != nil {
			logger.Error(
				"foundation boot failed",
				"error", err,
			)

			return err
		}

	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	return shutdownFoundation(f, logger)
}

func shutdownFoundation(f foundation.Foundation, logger common.Logger) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := f.Shutdown(ctx); err != nil {
		logger.Error(
			"foundation shutdown failed",
			"error", err,
		)

		return err
	}

	logger.Info("application shutdown completed")

	return nil
}
