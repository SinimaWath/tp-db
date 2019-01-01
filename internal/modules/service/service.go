package service

import (
	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) Clear(operations.ClearParams) middleware.Responder {
	err := database.Clear(self.db)
	if err != nil {
		return nil
	}

	return operations.NewClearOK()
}

func (self *ForumPgsql) Status(params operations.StatusParams) middleware.Responder {
	status := &models.Status{}
	err := database.Status(self.db, status)
	if err != nil {
		return nil
	}
	return operations.NewStatusOK().WithPayload(status)
}
