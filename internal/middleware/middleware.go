package middleware

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/rs/zerolog/hlog"
	"github.com/vrv501/simple-api/internal/constants"
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
				jsonBody, _ := json.Marshal(map[string]string{
					"msg": http.StatusText(http.StatusInternalServerError)})
				w.Write(jsonBody)
				w.Write([]byte("\n"))
			}
		}()
		h.ServeHTTP(w, r)
	})
}
