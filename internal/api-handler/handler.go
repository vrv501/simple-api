package apihandler

import (
	"context"

	"github.com/vrv501/simple-api/internal/constants"
	"github.com/vrv501/simple-api/internal/db"
)

type APIHandler struct {
	dbClient db.DBHandler
}

func NewAPIHandler(ctx context.Context) *APIHandler {
	return &APIHandler{
		dbClient: db.NewDBHandler(ctx),
	}
}

// Closes all clients associated with api handler
func (a *APIHandler) Close() {
	timedCtx, cancel := context.WithTimeout(context.Background(),
		constants.DefaultShutdownTimeout)
	defer cancel()
	_ = a.dbClient.Close(timedCtx)
}
