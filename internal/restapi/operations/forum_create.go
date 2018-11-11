// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ForumCreateHandlerFunc turns a function with the right signature into a forum create handler
type ForumCreateHandlerFunc func(ForumCreateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ForumCreateHandlerFunc) Handle(params ForumCreateParams) middleware.Responder {
	return fn(params)
}

// ForumCreateHandler interface for that can handle valid forum create params
type ForumCreateHandler interface {
	Handle(ForumCreateParams) middleware.Responder
}

// NewForumCreate creates a new http.Handler for the forum create operation
func NewForumCreate(ctx *middleware.Context, handler ForumCreateHandler) *ForumCreate {
	return &ForumCreate{Context: ctx, Handler: handler}
}

/*ForumCreate swagger:route POST /forum/create forumCreate

Создание форума

Создание нового форума.


*/
type ForumCreate struct {
	Context *middleware.Context
	Handler ForumCreateHandler
}

func (o *ForumCreate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewForumCreateParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
