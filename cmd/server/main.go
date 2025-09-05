package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
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
	ogenMw := ogenMiddleware.OapiRequestValidatorWithOptions(spec, &ogenMiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, // Update once authnz is implemented
		},
		ErrorHandlerWithOpts: func(_ context.Context, err error, w http.ResponseWriter,
			_ *http.Request, opts ogenMiddleware.ErrorHandlerOpts) {
			var reqErr *openapi3filter.RequestError
			if errors.As(err, &reqErr) {
				errorLines := strings.Split(reqErr.Error(), "\n")
				err = errors.New(errorLines[0])
			}

			w.Header().Del("Content-Length")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(opts.StatusCode)
			jsonBody, _ := json.Marshal(
				genRouter.Generic{
					Message: err.Error(),
				},
			)
			w.Write(jsonBody)
			w.Write([]byte("\n"))
		},
		SilenceServersWarning: true,
	})
	registerBodyEncoders()
	registerBodyDecoders()

	logger := configLogger()
	basePath, err := spec.Servers.BasePath()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get base path from OpenAPI spec")
	}

	port := getPort(logger)
	router := http.NewServeMux()
	router.HandleFunc(http.MethodGet+" /status",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	apiHandler := apihandler.NewAPIHandler(ctx)
	defer apiHandler.Close()

	routerWithCors := genRouter.HandlerWithOptions(
		genRouter.NewStrictHandler(apiHandler, nil),
		genRouter.StdHTTPServerOptions{
			BaseURL:    basePath,
			BaseRouter: router,
			Middlewares: []genRouter.MiddlewareFunc{
				/*
					Order matters here.
					Middleware is executed in order from reverse
				*/
				middleware.EntryAudit,
				hlog.RequestHandler("url"),
				hlog.RemoteAddrHandler("client_ip"),
				hlog.RequestIDHandler("request_id", "X-Request-ID"),
				hlog.AccessHandler(
					// The below function is a deferred call
					func(r *http.Request, status, _ int, duration time.Duration) {
						hlog.FromRequest(r).Info().
							Int("status", status).
							Str("latency", duration.String()).
							Msg("Exit Audit")
					},
				),
				ogenMw,
				middleware.PanicRecovery,
				hlog.NewHandler(logger),
			},
		},
	)
	routerWithCors = middleware.WithCORS(routerWithCors)

	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		Handler:      routerWithCors,
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
		logger.Info().Msgf("Started server on port %d", port)
		if errS := server.ListenAndServe(); errS != nil &&
			!errors.Is(errS, http.ErrServerClosed) {
			logger.Fatal().Err(errS).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()
	timedCtx, cancel := context.WithTimeout(context.Background(),
		constants.DefaultTimeout)
	defer cancel()
	if err = server.Shutdown(timedCtx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to shutdown server")
	}
}

func getPort(logger zerolog.Logger) int {
	portStr := os.Getenv(constants.ServerPort)
	if portStr == "" {
		return 8300
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

func registerBodyEncoders() {
	openapi3filter.RegisterBodyEncoder(
		"image/jpeg",
		func(body any) ([]byte, error) {
			if imgStr, ok := body.(string); ok {
				return []byte(imgStr), nil
			}
			return nil, errors.New("cannot encode image data: expected string (base64) or []byte")
		},
	)
}

func registerBodyDecoders() {
	openapi3filter.RegisterBodyDecoder("application/merge-patch+json", openapi3filter.JSONBodyDecoder)
	openapi3filter.RegisterBodyDecoder("image/jpeg", func(body io.Reader, _ http.Header, _ *openapi3.SchemaRef,
		_ openapi3filter.EncodingFn) (any, error) {
		buf := make([]byte, constants.MaxImgSize+1)
		n, err := io.ReadFull(body, buf)
		if err != nil && !errors.Is(err, io.EOF) &&
			!errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, err
		}
		if n == 0 || n > constants.MaxImgSize {
			return nil, errors.New("images should have min size 1B and max size 256KB")
		}
		return string(buf[:n]), nil
	})
}
