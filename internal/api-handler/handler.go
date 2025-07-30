package apihandler

import "github.com/vrv501/simple-api/internal/db"

type apiHandler struct {
	dbClient db.DBHandler
}

func NewAPIHandler() *apiHandler {
	return &apiHandler{
		dbClient: db.NewDBHandler(),
	}
}
