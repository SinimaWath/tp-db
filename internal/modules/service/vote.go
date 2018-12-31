package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) ThreadVote(params operations.ThreadVoteParams) middleware.Responder {
	log.Println("[INFO] ThreadVote")
	thread := &models.Thread{}
	err := database.VoteCreate(self.db, params.SlugOrID, thread, params.Vote)
	if err != nil {
		switch err {
		case database.ErrThreadNotFound:
			return operations.NewThreadVoteNotFound().WithPayload(&models.Error{})
		}
		log.Println("[ERROR] ThreadVote: " + err.Error())
		return nil
	}
	return operations.NewThreadVoteOK().WithPayload(thread)
}
