// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/SinimaWath/tp-db/internal/models"
)

// ForumGetOneOKCode is the HTTP code returned for type ForumGetOneOK
const ForumGetOneOKCode int = 200

/*ForumGetOneOK Информация о форуме.


swagger:response forumGetOneOK
*/
type ForumGetOneOK struct {

	/*
	  In: Body
	*/
	Payload *models.Forum `json:"body,omitempty"`
}

// NewForumGetOneOK creates ForumGetOneOK with default headers values
func NewForumGetOneOK() *ForumGetOneOK {

	return &ForumGetOneOK{}
}

// WithPayload adds the payload to the forum get one o k response
func (o *ForumGetOneOK) WithPayload(payload *models.Forum) *ForumGetOneOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the forum get one o k response
func (o *ForumGetOneOK) SetPayload(payload *models.Forum) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ForumGetOneOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ForumGetOneNotFoundCode is the HTTP code returned for type ForumGetOneNotFound
const ForumGetOneNotFoundCode int = 404

/*ForumGetOneNotFound Форум отсутсвует в системе.


swagger:response forumGetOneNotFound
*/
type ForumGetOneNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewForumGetOneNotFound creates ForumGetOneNotFound with default headers values
func NewForumGetOneNotFound() *ForumGetOneNotFound {

	return &ForumGetOneNotFound{}
}

// WithPayload adds the payload to the forum get one not found response
func (o *ForumGetOneNotFound) WithPayload(payload *models.Error) *ForumGetOneNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the forum get one not found response
func (o *ForumGetOneNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ForumGetOneNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}