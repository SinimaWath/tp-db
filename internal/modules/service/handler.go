package service

import (
	"errors"

	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

var (
	errNotFound            = errors.New("Not found")
	errInternalServer      = errors.New("Internal server error")
	errEmptySelect         = errors.New("Empty select")
	errUniqueViolation     = errors.New("Unique Violation")
	errForeignKeyViolation = errors.New("Foreign Key Violation")
)

type ForumHandler interface {
	ForumCreate(operations.ForumCreateParams) middleware.Responder
	UserCreate(operations.UserCreateParams) middleware.Responder
	UserGetOne(operations.UserGetOneParams) middleware.Responder
	UserUpdate(operations.UserUpdateParams) middleware.Responder
	ForumGetOne(operations.ForumGetOneParams) middleware.Responder
	ThreadCreate(operations.ThreadCreateParams) middleware.Responder
	ForumGetThreads(operations.ForumGetThreadsParams) middleware.Responder
	Clear(operations.ClearParams) middleware.Responder
}
