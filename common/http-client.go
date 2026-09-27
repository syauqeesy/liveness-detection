package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CommonHttpClient interface {
	PostJSON(
		ctx context.Context,
		url string,
		headers map[string]string,
		payload any,
		response any,
	) error
}

type commonHttpClient struct {
	client *http.Client
	logger Logger
}

func NewHttpClient(logger Logger) *commonHttpClient {
	return &commonHttpClient{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

func (h *commonHttpClient) PostJSON(
	ctx context.Context,
	url string,
	headers map[string]string,
	payload any,
	response any,
) error {
	started := time.Now()

	body, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error(
			"failed to marshal HTTP request",
			"error", err,
			"url", url,
		)

		return fmt.Errorf("marshal HTTP request: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		h.logger.Error(
			"failed to create HTTP request",
			"error", err,
			"url", url,
		)

		return fmt.Errorf("create HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	h.logger.Debug(
		"HTTP request started",
		"method", request.Method,
		"url", url,
	)

	result, err := h.client.Do(request)
	if err != nil {
		h.logger.Error(
			"HTTP request failed",
			"error", err,
			"method", request.Method,
			"url", url,
			"duration_ms", time.Since(started).Seconds()*1000,
		)

		return fmt.Errorf("execute HTTP request: %w", err)
	}

	defer result.Body.Close()

	durationMs := time.Since(started).Seconds() * 1000

	if result.StatusCode < 200 || result.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(result.Body)

		h.logger.Warn(
			"HTTP request rejected",
			"method", request.Method,
			"url", url,
			"status", result.Status,
			"duration_ms", durationMs,
		)

		return fmt.Errorf(
			"HTTP request failed: status=%s body=%s",
			result.Status,
			string(responseBody),
		)
	}

	if response != nil {
		if err := json.NewDecoder(result.Body).Decode(response); err != nil {
			h.logger.Error(
				"failed to decode HTTP response",
				"error", err,
				"method", request.Method,
				"url", url,
				"status", result.Status,
				"duration_ms", durationMs,
			)

			return fmt.Errorf("decode HTTP response: %w", err)
		}
	}

	h.logger.Info(
		"HTTP request completed",
		"method", request.Method,
		"url", url,
		"status", result.Status,
		"duration_ms", durationMs,
	)

	return nil
}
