package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	ogenMiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	apihandler "github.com/vrv501/simple-api/internal/api-handler"
	"github.com/vrv501/simple-api/internal/constants"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
	"github.com/vrv501/simple-api/internal/middleware"
)

func main() {
	ctx := signals.SetupSignalHandler()
	spec, _ := genRouter.GetSwagger()
	spec.Servers = nil
	ogenMw := ogenMiddleware.OapiRequestValidatorWithOptions(spec, &ogenMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(_ context.Context, _ *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
	})

	logger := configLogger()
	port := getPort(logger)
	router := http.NewServeMux()
	router.HandleFunc(http.MethodGet+" "+constants.StatusPath,
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			byteArr, _ := json.Marshal(map[string]string{"status": "healthy"})
			w.Write(byteArr)
			w.Write([]byte("\n"))
		},
	)

	apiHandler := apihandler.NewAPIHandler(ctx)
	defer apiHandler.Close()

	server := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
		Handler: genRouter.HandlerWithOptions(
			genRouter.NewStrictHandler(apiHandler, nil),
			genRouter.StdHTTPServerOptions{
				BaseRouter: router,
				Middlewares: []genRouter.MiddlewareFunc{
					/*
						Order matters here.
						Middleware is executed in order from reverse
					*/
					middleware.EntryAudit,
					hlog.RequestHandler(constants.LogFieldMethodAndURL),
					hlog.RemoteAddrHandler(constants.LogFieldClientIP),
					hlog.RequestIDHandler(constants.LogFieldRequestID, constants.HeaderRequestID),
					hlog.AccessHandler(
						// The below function is a deferred call
						func(r *http.Request, status, _ int, duration time.Duration) {
							hlog.FromRequest(r).Info().
								Int(constants.LogFieldStatus, status).
								Str(constants.LogFieldLatency, duration.String()).
								Msg("Exit Audit")
						},
					),
					ogenMw,
					middleware.PanicRecovery,
					hlog.NewHandler(logger),
				},
			},
		),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	go func() {
		// We need not monitor this go-routine
		// When [http.Server] Shutdown is called, ListenAndServe() immediately returns
		// [http.ErrServerClosed]. Shutdown() takes some-time to cleanup & gracefully close the server
		// That being said, if server initiate itself failed with some error
		// log.Fatal() will call os.Exit(1) which halts the entire program
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()
	timedCtx, cancel := context.WithTimeout(context.Background(),
		constants.DefaultShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(timedCtx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to shutdown server")
	}
}

func getPort(logger zerolog.Logger) int {
	portStr := os.Getenv(constants.ServerPort)
	if portStr == "" {
		return constants.DefaultServerPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Invalid %s", constants.ServerPort)
	}
	return port
}

func configLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	// To disable logging entirely, pass [zerolog.Disabled]
	logLevel, err := zerolog.ParseLevel(os.Getenv(constants.LogLevel))
	if err != nil {
		logger.Fatal().Err(err).Msgf("Invalid %s", constants.LogLevel)
	}
	if logLevel.String() == "" {
		logLevel = zerolog.InfoLevel
	}

	return logger.Level(logLevel)
}
