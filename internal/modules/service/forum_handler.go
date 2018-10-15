package service

import (
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type ForumHandler interface {
	ForumCreate(operations.ForumCreateParams) middleware.Responder
}
