package apihandler

import "github.com/vrv501/simple-api/internal/db"

type ApiHandler struct {
	dbClient db.DBHandler
}

func NewAPIHandler() *ApiHandler {
	return &ApiHandler{
		//dbClient: db.NewDBHandler(),
	}
}
