package ceng_pubsub

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentCash PaymentMethod = "cash"
	PaymentCard PaymentMethod = "card"
)

type PrinterEventEntity struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Active    *bool     `json:"active"`
	Inside    *bool     `json:"inside"`
	Outside   *bool     `json:"outside"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MenuCategoryEventEntity struct {
	ID               uuid.UUID  `json:"id"`
	Title            string     `json:"title"`
	Position         int64      `json:"position"`
	Active           *bool      `json:"active"`
	Inside           *bool      `json:"inside"`
	Outside          *bool      `json:"outside"`
	PrinterInsideID  *uuid.UUID `json:"printerInsideId"`
	PrinterOutsideID *uuid.UUID `json:"printerOutsideId"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type MenuItemEventEntity struct {
	ID             uuid.UUID `json:"id"`
	MenuCategoryID uuid.UUID `json:"menuCategoryId"`
	Title          string    `json:"title"`
	Position       int64     `json:"position"`
	Active         *bool     `json:"active"`
	Inside         *bool     `json:"inside"`
	Outside        *bool     `json:"outside"`
	Price          int64     `json:"price"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type MenuOptionEventEntity struct {
	ID         uuid.UUID `json:"id"`
	MenuItemID uuid.UUID `json:"menuItemId"`
	Title      string    `json:"title"`
	Position   int64     `json:"position"`
	Active     *bool     `json:"active"`
	Inside     *bool     `json:"inside"`
	Outside    *bool     `json:"outside"`
	Price      int64     `json:"price"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TableEventEntity struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userId"`
	Name          string    `json:"name"`
	Close         *bool     `json:"close"`
	Inside        *bool     `json:"inside"`
	PaymentMethod *string   `json:"paymentMethod"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type OrderEventEntity struct {
	ID        uuid.UUID `json:"id"`
	TableID   uuid.UUID `json:"tableId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CourseEventEntity struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"orderId"`
	Position  int64     `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CourseSelectionEventEntity struct {
	ID           uuid.UUID  `json:"id"`
	CourseID     uuid.UUID  `json:"courseId"`
	MenuItemID   uuid.UUID  `json:"menuItemId"`
	MenuOptionID *uuid.UUID `json:"menuOptionId"`
	Quantity     int64      `json:"quantity"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
