// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostGetOneHandlerFunc turns a function with the right signature into a post get one handler
type PostGetOneHandlerFunc func(PostGetOneParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostGetOneHandlerFunc) Handle(params PostGetOneParams) middleware.Responder {
	return fn(params)
}

// PostGetOneHandler interface for that can handle valid post get one params
type PostGetOneHandler interface {
	Handle(PostGetOneParams) middleware.Responder
}

// NewPostGetOne creates a new http.Handler for the post get one operation
func NewPostGetOne(ctx *middleware.Context, handler PostGetOneHandler) *PostGetOne {
	return &PostGetOne{Context: ctx, Handler: handler}
}

/*PostGetOne swagger:route GET /post/{id}/details postGetOne

Получение информации о ветке обсуждения

Получение информации о ветке обсуждения по его имени.


*/
type PostGetOne struct {
	Context *middleware.Context
	Handler PostGetOneHandler
}

func (o *PostGetOne) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostGetOneParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
