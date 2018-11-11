// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Pet pet
// swagger:model Pet
type Pet struct {

	// category
	Category *Category `json:"category,omitempty"`

	// id
	ID int64 `json:"id,omitempty" gorm:"primary_key" query:"filter,sort"`

	// name
	// Required: true
	Name *string `json:"name" query:"filter,sort"`

	// photo urls
	// Required: true
	PhotoUrls []string `json:"photoUrls" gorm:"-"`

	// pet status in the store
	// Enum: [available pending sold]
	Status string `json:"status,omitempty" query:"filter,sort"`

	// tags
	Tags []*Tag `json:"tags"`
}

// Validate validates this pet
func (m *Pet) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCategory(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePhotoUrls(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Pet) validateCategory(formats strfmt.Registry) error {

	if swag.IsZero(m.Category) { // not required
		return nil
	}

	if m.Category != nil {
		if err := m.Category.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("category")
			}
			return err
		}
	}

	return nil
}

func (m *Pet) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Pet) validatePhotoUrls(formats strfmt.Registry) error {

	if err := validate.Required("photoUrls", "body", m.PhotoUrls); err != nil {
		return err
	}

	return nil
}

var petTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["available","pending","sold"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		petTypeStatusPropEnum = append(petTypeStatusPropEnum, v)
	}
}

const (

	// PetStatusAvailable captures enum value "available"
	PetStatusAvailable string = "available"

	// PetStatusPending captures enum value "pending"
	PetStatusPending string = "pending"

	// PetStatusSold captures enum value "sold"
	PetStatusSold string = "sold"
)

// prop value enum
func (m *Pet) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, petTypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Pet) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

func (m *Pet) validateTags(formats strfmt.Registry) error {

	if swag.IsZero(m.Tags) { // not required
		return nil
	}

	for i := 0; i < len(m.Tags); i++ {
		if swag.IsZero(m.Tags[i]) { // not required
			continue
		}

		if m.Tags[i] != nil {
			if err := m.Tags[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Pet) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Pet) UnmarshalBinary(b []byte) error {
	var res Pet
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
