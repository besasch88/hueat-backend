package menuCategory

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type createMenuCategoryInputDto struct {
	PrinterInsideID  *string `json:"printerInsideId"`
	PrinterOutsideID *string `json:"printerOutsideId"`
	Title            string  `json:"title"`
}

func (r createMenuCategoryInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.PrinterInsideID, is.UUID),
		validation.Field(&r.PrinterOutsideID, is.UUID),
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
	)
}

type getMenuCategoryInputDto struct {
	ID string `uri:"menuCategoryId"`
}

func (r getMenuCategoryInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type updateMenuCategoryInputDto struct {
	ID               string  `uri:"menuCategoryId"`
	PrinterInsideID  *string `json:"printerInsideId"`
	PrinterOutsideID *string `json:"printerOutsideId"`
	Title            *string `json:"title"`
	Position         *int64  `json:"position"`
	Active           *bool   `json:"active"`
	Inside           *bool   `json:"inside"`
	Outside          *bool   `json:"outside"`
}

func (r updateMenuCategoryInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.PrinterInsideID, is.UUID),
		validation.Field(&r.PrinterOutsideID, is.UUID),
		validation.Field(&r.Title, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Position, validation.NilOrNotEmpty, validation.Min(1)),
		validation.Field(&r.Active, validation.In(true, false)),
		validation.Field(&r.Inside, validation.In(true, false)),
		validation.Field(&r.Outside, validation.In(true, false)),
	)
}

type deleteMenuCategoryInputDto struct {
	ID string `uri:"menuCategoryId"`
}

func (r deleteMenuCategoryInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
