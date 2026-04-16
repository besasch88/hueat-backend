package menu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type getMenuInputDto struct {
	TableID string `uri:"tableId"`
}

func (r getMenuInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TableID, validation.Required, is.UUID),
	)
}
