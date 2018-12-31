package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) Clear(operations.ClearParams) middleware.Responder {
	log.Println("[INFO] Clear")
	err := database.Clear(self.db)
	if err != nil {
		log.Println("[ERROR] Clear: " + err.Error())
		return nil
	}

	return operations.NewClearOK()
}

func (self *ForumPgsql) Status(params operations.StatusParams) middleware.Responder {
	log.Println("[INFO] Status")
	status := &models.Status{}
	err := database.Status(self.db, status)
	if err != nil {
		log.Println("[ERROR] Status: " + err.Error())
		return nil
	}
	return operations.NewStatusOK().WithPayload(status)
}
