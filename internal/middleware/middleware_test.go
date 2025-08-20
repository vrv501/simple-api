package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vrv501/simple-api/internal/constants"
)

func TestEntryAudit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		wantStatusCode int
	}{
		{
			name:           "EntryAudit",
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler := EntryAudit(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
			if rr.Code != tt.wantStatusCode {
				t.Errorf("EntryAudit() status code = %v, want %v", rr.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	type args struct {
		h http.Handler
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "PanicRecovery",
			args: args{
				h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					panic("test panic")
				}),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler := PanicRecovery(tt.args.h)
			handler.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
			if rr.Code != tt.wantStatusCode {
				t.Errorf("PanicRecovery() = %v, want %v", rr.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWithCORS(t *testing.T) {
	t.Setenv(constants.AllowedOrigins, "http://localhost:8080")

	type args struct {
		req *http.Request
	}
	tests := []struct {
		name           string
		args           args
		prepare        func(*args)
		wantStatusCode int
	}{
		{
			name: "WithCORS",
			args: args{
				req: httptest.NewRequest("GET", "/", nil),
			},
			prepare: func(args *args) {
				args.req.Header.Set("Origin", "http://localhost:8080")
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "WithCORS Options",
			args: args{
				req: httptest.NewRequest("OPTIONS", "/", nil),
			},
			wantStatusCode: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			if tt.prepare != nil {
				tt.prepare(&tt.args)
			}
			handler := WithCORS(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			handler.ServeHTTP(rr, tt.args.req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("WithCORS() = %v, want %v", rr.Code, tt.wantStatusCode)
			}
		})
	}
}
