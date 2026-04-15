package printer

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type createPrinterInputDto struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func (r createPrinterInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Url, validation.Required, validation.Length(1, 255)),
	)
}

type getPrinterInputDto struct {
	ID string `uri:"printerId"`
}

func (r getPrinterInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type updatePrinterInputDto struct {
	ID     string  `uri:"printerId"`
	Title  *string `json:"title"`
	Url    *string `json:"url"`
	Active *bool   `json:"active"`
}

func (r updatePrinterInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Url, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Active, validation.In(true, false)),
	)
}

type deletePrinterInputDto struct {
	ID string `uri:"printerId"`
}

func (r deletePrinterInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
