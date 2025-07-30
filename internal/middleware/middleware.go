package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vrv501/simple-api/internal/constants"
)

func AddRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.URL.Path != "/status" {
			requestID := r.Header.Get(constants.HeaderRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			} else {
				err := uuid.Validate(requestID)
				if err != nil {
					requestID = uuid.New().String()
				}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.RequestIDKey, requestID)
			ctx = context.WithValue(ctx, constants.ClientIPAddr, r.RemoteAddr)
			ctx = context.WithValue(ctx, constants.HTTPMethod, r.Method)
			ctx = context.WithValue(ctx, constants.URLPath, r.URL.Path)
			r = r.WithContext(ctx)

			r.Header.Set(constants.HeaderRequestID, requestID)
			w.Header().Set(constants.HeaderRequestID, requestID)
		}
		h.ServeHTTP(w, r)
	})
}

func PanicRecovery(h http.Handler) http.Handler {
	jsonBody, _ := json.Marshal(`{"msg": "Internal Server Error"}`)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]
				log.Ctx(r.Context()).Error().
					Interface(constants.LogFieldPanic, err).
					Str(constants.LogFieldStackTrace, string(stack)).
					Msg("Recovered from panic")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
