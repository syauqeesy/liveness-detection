package handler

import (
	"net/http"

	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
	"github.com/syauqeesy/liveness-detection/service"
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
