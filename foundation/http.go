package foundation

import (
	"context"
	"net/http"

	"github.com/syauqeesy/liveness-detection/common"
	"github.com/syauqeesy/liveness-detection/configuration"
	"github.com/syauqeesy/liveness-detection/handler"
	"github.com/syauqeesy/liveness-detection/middleware"
	grpc_outbound "github.com/syauqeesy/liveness-detection/outbound/grpc"
	"github.com/syauqeesy/liveness-detection/service"
)

type httpFoundation struct {
	configuration       *configuration.Configuration
	mux                 *http.ServeMux
	server              *http.Server
	service             *service.Service
	handler             *handler.Handler
	logger              common.Logger
	http                common.CommonHttp
	grpcOutboundService *grpc_outbound.GRPCOutboundService
}

func (f *httpFoundation) Setup() error {
	f.mux = http.NewServeMux()

	f.http = common.NewHttp(f.logger)

	f.grpcOutboundService = grpc_outbound.New(f.configuration)

	f.service = service.NewService(f.configuration, f.grpcOutboundService)

	f.handler = handler.NewHandler(f.mux, f.configuration, f.service, f.http)

	f.logger = common.NewLogger(f.configuration.Application.Service, f.configuration.Application.Environment)

	f.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f.http.ErrorHandler(w, common.CreateException(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)), nil)
	})

	f.server = &http.Server{
		Addr:    f.configuration.Http.Port,
		Handler: middleware.Cors(f.configuration)(middleware.Logger(f.logger)(f.mux)),
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
