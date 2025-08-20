package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/rs/zerolog"
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
					Interface(zerolog.ErrorFieldName, err).
					Str("stack_trace", string(stack)).
					Msg("Recovered from panic")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				jsonBody, _ := json.Marshal(
					genRouter.Generic{
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
	allowedOriginList := []string{"http://localhost:8080"}
	if allowedOrigins := os.Getenv(constants.AllowedOrigins); allowedOrigins != "" {
		for allowedOrigin := range strings.SplitSeq(allowedOrigins, ",") {
			allowedOriginList = append(allowedOriginList, strings.TrimSpace(allowedOrigin))
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "X-Request-Id, X-Next-Cursor")
		if allowedOrigin := r.Header.Get("Origin"); slices.Contains(allowedOriginList, allowedOrigin) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
