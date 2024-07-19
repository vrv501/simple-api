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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/vrv501/go-template/internal/constants"
	"github.com/vrv501/go-template/internal/generated"
)

func main() {
	ctx := signals.SetupSignalHandler()

	configLogger()
	port := getPort()

	router := http.NewServeMux()
	router.HandleFunc("GET /status",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			byteArr, _ := json.Marshal(map[string]string{"status": "healthy"})
			w.Write(byteArr)
		},
	)
	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		Handler:      generated.HandlerFromMux(_, router),
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
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()
	timedCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := server.Shutdown(timedCtx); err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}
}

func getPort() int {
	portStr := os.Getenv(constants.ServerPort)
	if portStr == "" {
		return constants.DefaultServerPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal().Err(err).Msgf("Invalid %s", constants.ServerPort)
	}
	return port
}

func configLogger() {
	log.Logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	logLevel, err := zerolog.ParseLevel(os.Getenv(constants.LogLevel))
	if err != nil {
		log.Fatal().Err(err).Msgf("Invalid %s", constants.LogLevel)
	}
	if logLevel.String() == "" {
		logLevel = zerolog.InfoLevel
	}

	// To disable logging entirely, pass [zerolog.Disabled]
	zerolog.SetGlobalLevel(logLevel)
}
