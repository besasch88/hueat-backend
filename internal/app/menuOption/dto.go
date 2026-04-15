package menuOption

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type listMenuOptionsInputDto struct {
	MenuItemId string `uri:"menuItemId"`
}

func (r listMenuOptionsInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuItemId, validation.Required, is.UUID),
	)
}

type createMenuOptionInputDto struct {
	MenuItemId   string `uri:"menuItemId"`
	Title        string `json:"title"`
	TitleDisplay string `json:"TitleDisplay"`
	Price        int64  `json:"price"`
}

func (r createMenuOptionInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuItemId, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.TitleDisplay, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Price, validation.Required, validation.Min(1), validation.Max(10000)),
	)
}

type getMenuOptionInputDto struct {
	ID string `uri:"menuOptionId"`
}

func (r getMenuOptionInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type updateMenuOptionInputDto struct {
	ID           string  `uri:"menuOptionId"`
	Title        *string `json:"title"`
	TitleDisplay *string `json:"TitleDisplay"`
	Position     *int64  `json:"position"`
	Active       *bool   `json:"active"`
	Inside       *bool   `json:"inside"`
	Outside      *bool   `json:"outside"`
	Price        *int64  `json:"price"`
}

func (r updateMenuOptionInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.TitleDisplay, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Position, validation.NilOrNotEmpty, validation.Min(1)),
		validation.Field(&r.Active, validation.In(true, false)),
		validation.Field(&r.Inside, validation.In(true, false)),
		validation.Field(&r.Outside, validation.In(true, false)),
		validation.Field(&r.Price, validation.Min(1), validation.Max(10000)),
	)
}

type deleteMenuOptionInputDto struct {
	ID string `uri:"menuOptionId"`
}

func (r deleteMenuOptionInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
