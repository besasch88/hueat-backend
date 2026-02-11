package order

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
)

type orderEntity ceng_pubsub.OrderEventEntity
type courseEntity ceng_pubsub.CourseEventEntity
type courseSelectionEntity ceng_pubsub.CourseSelectionEventEntity

type courseEntityWithChilds struct {
	courseEntity
	Items []courseSelectionEntity `json:"items"`
}

type orderEntityWithChilds struct {
	orderEntity
	Courses []courseEntityWithChilds `json:"courses"`
}

type orderDetailEntity struct {
	Username        string
	TableName       string
	TableCreatedAt  time.Time
	PrinterId       string
	PrinterTitle    string
	PrinterURL      string
	CourseID        string
	CourseNumber    int64
	MenuItemTitle   string
	MenuItemPrice   int64
	MenuOptionTitle *string
	MenuOptionPrice *int64
	Quantity        int64
}

type paymentDetailEntity struct {
	Username       string
	TableName      string
	TableCreatedAt time.Time
	TablePayment   string
	PrinterId      string
	PrinterTitle   string
	PrinterURL     string
	PriceTotal     int64
}
