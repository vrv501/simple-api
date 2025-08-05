package apihandler

import (
	"context"

	"github.com/vrv501/simple-api/internal/constants"
	"github.com/vrv501/simple-api/internal/db"
)

type ApiHandler struct {
	dbClient db.DBHandler
}

func NewAPIHandler(ctx context.Context) *ApiHandler {
	return &ApiHandler{
		dbClient: db.NewDBHandler(ctx),
	}
}

// Closes all clients associated with api handler
func (a *ApiHandler) Close() {
	timedCtx, cancel := context.WithTimeout(context.Background(),
		constants.DefaultShutdownTimeout)
	defer cancel()
	_ = a.dbClient.Close(timedCtx)
}
