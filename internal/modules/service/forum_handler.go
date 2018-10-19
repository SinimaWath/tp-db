package service

import (
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type ForumHandler interface {
	ForumCreate(operations.ForumCreateParams) middleware.Responder
	UserCreate(operations.UserCreateParams) middleware.Responder
	UserGetOne(operations.UserGetOneParams) middleware.Responder
	UserUpdate(operations.UserUpdateParams) middleware.Responder
	ForumGetOne(operations.ForumGetOneParams) middleware.Responder
	ThreadCreate(operations.ThreadCreateParams) middleware.Responder
	ForumGetThreads(operations.ForumGetThreadsParams) middleware.Responder
}
