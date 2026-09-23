package handler

import (
	"net/http"

	"github.com/liveness-detection/common"
	"github.com/liveness-detection/configuration"
	"github.com/liveness-detection/service"
)

type handler struct {
	Service       *service.Service
	Configuration *configuration.Configuration
	CommonHttp    common.CommonHttp
}

type Handler struct {
	Inference *inferenceHandler
}

func NewHandler(mux *http.ServeMux, configuration *configuration.Configuration, service *service.Service, commonHttp common.CommonHttp) *Handler {
	handler := &handler{
		Configuration: configuration,
		Service:       service,
		CommonHttp:    commonHttp,
	}

	h := &Handler{
		Inference: (*inferenceHandler)(handler),
	}

	mux.HandleFunc("POST /inference", h.Inference.Execute)

	return h
}
