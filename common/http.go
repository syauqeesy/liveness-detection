package common

import (
	"encoding/json"
	"net/http"
)

type HttpJsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type CommonHttp interface {
	ErrorHandler(w http.ResponseWriter, exception error, payload any)
	WriteResponse(w http.ResponseWriter, status int, message string, payload any)
	Redirect(w http.ResponseWriter, request *http.Request, url string)
}

type commonHttp struct {
	logger Logger
}

func NewHttp(logger Logger) *commonHttp {
	return &commonHttp{
		logger: logger,
	}
}

func (h *commonHttp) ErrorHandler(w http.ResponseWriter, exception error, payload any) {
	w.Header().Set("Content-Type", "application/json")

	message := http.StatusText(http.StatusInternalServerError)
	status := http.StatusInternalServerError

	switch convertedException := exception.(type) {
	case *ApplicationException:
		status = convertedException.HttpStatusCode
		message = convertedException.Error()

		h.logger.Info("http request rejected", "error", convertedException)
	default:
		status = http.StatusInternalServerError
		message = http.StatusText(http.StatusInternalServerError)

		h.logger.Error(
			"unhandled application exception",
			"error", exception,
		)
	}

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&HttpJsonResponse{
		Message: message,
		Data:    payload,
	}); err != nil {
		h.logger.Error(
			"failed to encode http error response",
			"error", err,
		)
	}
}

func (h *commonHttp) WriteResponse(w http.ResponseWriter, status int, message string, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&HttpJsonResponse{
		Message: message,
		Data:    payload,
	}); err != nil {
		h.logger.Error(
			"failed to encode http response",
			"error", err,
		)
	}
}

func (h *commonHttp) Redirect(w http.ResponseWriter, request *http.Request, url string) {
	http.Redirect(w, request, url, http.StatusFound)
}
