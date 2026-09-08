package handler

import (
	"net/http"

	"ahmadsyauqi.dev/projects/liveness-detection/configuration"
	"ahmadsyauqi.dev/projects/liveness-detection/service"
)

type handler struct {
	Service       *service.Service
	Configuration *configuration.Configuration
}

type Handler struct {
	Inference *inferenceHandler
}

func NewHandler(mux *http.ServeMux, configuration *configuration.Configuration, service *service.Service) *Handler {
	handler := &handler{
		Configuration: configuration,
		Service:       service,
	}

	h := &Handler{
		Inference: (*inferenceHandler)(handler),
	}

	return h
}
