package service

import (
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
)

type ForumPgsql struct {
	db *sqlx.DB
}

func (ForumPgsql) ForumCreate(operations.ForumCreateParams) middleware.Responder {

}
