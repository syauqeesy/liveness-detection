package foundation

import (
	"context"
	"net/http"

	"ahmadsyauqi.dev/projects/liveness-detection/common"
	"ahmadsyauqi.dev/projects/liveness-detection/configuration"
	"ahmadsyauqi.dev/projects/liveness-detection/handler"
	"ahmadsyauqi.dev/projects/liveness-detection/middleware"
	"ahmadsyauqi.dev/projects/liveness-detection/service"
)

type httpFoundation struct {
	configuration *configuration.Configuration
	mux           *http.ServeMux
	server        *http.Server
	service       *service.Service
	handler       *handler.Handler
	logger        common.Logger
	http          common.CommonHttp
}

func (f *httpFoundation) Setup() error {
	f.mux = http.NewServeMux()

	f.service = service.NewService(f.configuration)

	f.handler = handler.NewHandler(f.mux, f.configuration, f.service)

	f.logger = common.NewLogger(f.configuration.Application.Service, f.configuration.Application.Environment)

	f.http = common.NewHttp(f.logger)

	f.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f.http.ErrorHandler(w, common.CreateException(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)), nil)
	})

	f.server = &http.Server{
		Addr:    f.configuration.Http.Port,
		Handler: middleware.Logger(f.logger)(f.mux),
	}

	return nil
}

func (f *httpFoundation) Boot() error {
	f.logger.Info("http server started")

	err := f.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (f *httpFoundation) Shutdown(ctx context.Context) error {
	f.logger.Info("shutting down http server")

	err := f.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	f.logger.Info("http server exited")

	return nil
}
