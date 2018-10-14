package service

import (
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

type HelloHandler interface {
	AddMulti(params operations.AddMultiParams) middleware.Responder
	DestroyOne(params operations.DestroyOneParams) middleware.Responder
	Find(params operations.FindParams) middleware.Responder
	GetOne(params operations.GetOneParams) middleware.Responder
	UpdateOne(params operations.UpdateOneParams) middleware.Responder
}
