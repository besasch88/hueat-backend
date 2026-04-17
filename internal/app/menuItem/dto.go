package menuItem

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type listMenuItemsInputDto struct {
	MenuCategoryId string `uri:"menuCategoryId"`
}

func (r listMenuItemsInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuCategoryId, validation.Required, is.UUID),
	)
}

type createMenuItemInputDto struct {
	MenuCategoryId   string  `uri:"menuCategoryId"`
	Title            string  `json:"title"`
	TitleDisplay     string  `json:"TitleDisplay"`
	Price            int64   `json:"price"`
	PrinterInsideID  *string `json:"printerInsideId"`
	PrinterOutsideID *string `json:"printerOutsideId"`
}

func (r createMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuCategoryId, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.TitleDisplay, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Price, validation.Required, validation.Min(1), validation.Max(10000)),
		validation.Field(&r.PrinterInsideID, is.UUID),
		validation.Field(&r.PrinterOutsideID, is.UUID),
	)
}

type createCustomMenuItemInputDto struct {
	MenuCategoryId string `uri:"menuCategoryId"`
	TableId        string `uri:"tableId"`
	Title          string `json:"title"`
	TitleDisplay   string `json:"TitleDisplay"`
	Price          int64  `json:"price"`
	PrinterID      string `json:"printerId"`
}

func (r createCustomMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuCategoryId, validation.Required, is.UUID),
		validation.Field(&r.TableId, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.TitleDisplay, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Price, validation.Required, validation.Min(1), validation.Max(10000)),
		validation.Field(&r.PrinterID, validation.Required, is.UUID),
	)
}

type getMenuItemInputDto struct {
	ID string `uri:"menuItemId"`
}

func (r getMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type updateMenuItemInputDto struct {
	ID               string  `uri:"menuItemId"`
	Title            *string `json:"title"`
	TitleDisplay     *string `json:"TitleDisplay"`
	Position         *int64  `json:"position"`
	Active           *bool   `json:"active"`
	Inside           *bool   `json:"inside"`
	Outside          *bool   `json:"outside"`
	Price            *int64  `json:"price"`
	PrinterInsideID  *string `json:"printerInsideId"`
	PrinterOutsideID *string `json:"printerOutsideId"`
}

func (r updateMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.TitleDisplay, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Position, validation.NilOrNotEmpty, validation.Min(1)),
		validation.Field(&r.Active, validation.In(true, false)),
		validation.Field(&r.Inside, validation.In(true, false)),
		validation.Field(&r.Outside, validation.In(true, false)),
		validation.Field(&r.Price, validation.Min(1), validation.Max(10000)),
		validation.Field(&r.PrinterInsideID, is.UUID),
		validation.Field(&r.PrinterOutsideID, is.UUID),
	)
}

type deleteMenuItemInputDto struct {
	ID string `uri:"menuItemId"`
}

func (r deleteMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
