package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (pg ForumPgsql) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	queryInsert := `INSERT INTO forum (user_nick, slug, title) VALUES ($1, $2, $3)`

	_, err := pg.db.Exec(queryInsert, params.Forum.User, params.Forum.Slug, params.Forum.Title)
	if err != nil {
		log.Println(err)
		return nil
	}

	return operations.NewForumCreateCreated().WithPayload(params.Forum)
}
