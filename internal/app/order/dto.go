package order

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type getOrderInputDto struct {
	TableID string `uri:"tableId"`
}

func (r getOrderInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TableID, validation.Required, is.UUID),
	)
}

type updateCourseSelectionInputDto struct {
	MenuItemID   string  `json:"menuItemId"`
	MenuOptionID *string `json:"menuOptionId"`
	Quantity     int64   `json:"quantity"`
	Note         *string `json:"note"`
}

func (r updateCourseSelectionInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuItemID, validation.Required, is.UUID),
		validation.Field(&r.MenuOptionID, is.UUID),
		validation.Field(&r.Quantity, validation.Min(1)),
		validation.Field(&r.Note, validation.NilOrNotEmpty, validation.Length(1, 512)),
	)
}

type updateCourseInputDto struct {
	ID    string                          `json:"id"`
	Items []updateCourseSelectionInputDto `json:"items"`
}

func (r updateCourseInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Items, validation.Each(validation.By(func(value interface{}) error {
			v := value.(updateCourseSelectionInputDto)
			return v.validate()
		}))),
	)
}

type updateOrderInputDto struct {
	TableID string                 `uri:"tableId"`
	Courses []updateCourseInputDto `json:"courses"`
}

func (r updateOrderInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TableID, validation.Required, is.UUID),
		validation.Field(&r.Courses, validation.Required, validation.Length(1, 0), validation.Each(validation.By(func(value interface{}) error {
			v := value.(updateCourseInputDto)
			return v.validate()
		}))),
	)
}

type printOrderInputDto struct {
	TableID  string  `uri:"tableId"`
	CourseID *string `json:"courseId"`
	Target   string  `json:"target"`
}

func (r printOrderInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TableID, validation.Required, is.UUID),
		validation.Field(&r.CourseID, is.UUID, validation.NilOrNotEmpty, validation.When(r.Target == "course", validation.Required)),
		validation.Field(&r.Target, validation.Required, validation.In("order", "course", "bill", "payment")),
	)
}
