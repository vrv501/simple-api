package middleware

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/rs/zerolog/hlog"

	"github.com/vrv501/simple-api/internal/constants"
	genRouter "github.com/vrv501/simple-api/internal/generated/router"
)

func EntryAudit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hlog.FromRequest(r).Info().Msg("Entry Audit")
		h.ServeHTTP(w, r)
	})
}

func PanicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]
				hlog.FromRequest(r).Error().
					Interface(constants.LogFieldPanic, err).
					Str(constants.LogFieldStackTrace, string(stack)).
					Msg("Recovered from panic")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				jsonBody, _ := json.Marshal(
					genRouter.ApiResponse{
						Message: http.StatusText(http.StatusInternalServerError),
					},
				)
				w.Write(jsonBody)
				w.Write([]byte("\n"))
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func WithCORS(h http.Handler) http.Handler {
	allowedOrigin := "http://localhost:8080"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Header.Get("Origin") == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
