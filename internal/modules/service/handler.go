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
	ForumGetThreads(operations.ForumGetThreadsParams) middleware.Responder
	ForumGetUsers(operations.ForumGetUsersParams) middleware.Responder
	ForumGetOne(operations.ForumGetOneParams) middleware.Responder

	UserCreate(operations.UserCreateParams) middleware.Responder
	UserGetOne(operations.UserGetOneParams) middleware.Responder
	UserUpdate(operations.UserUpdateParams) middleware.Responder

	PostsCreate(operations.PostsCreateParams) middleware.Responder
	PostGetOne(operations.PostGetOneParams) middleware.Responder
	PostUpdate(operations.PostUpdateParams) middleware.Responder

	ThreadCreate(operations.ThreadCreateParams) middleware.Responder
	ThreadGetOne(operations.ThreadGetOneParams) middleware.Responder
	ThreadUpdate(operations.ThreadUpdateParams) middleware.Responder
	ThreadVote(operations.ThreadVoteParams) middleware.Responder
	ThreadGetPosts(operations.ThreadGetPostsParams) middleware.Responder

	Clear(operations.ClearParams) middleware.Responder
	Status(operations.StatusParams) middleware.Responder
}
