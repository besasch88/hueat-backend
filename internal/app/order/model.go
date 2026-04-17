package order

import (
	"time"

	"github.com/google/uuid"
)

type tableModel struct {
	ID     uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	UserID uuid.UUID `gorm:"column:user_id;type:varchar(36)"`
}

func (m tableModel) TableName() string {
	return "hueat_table"
}

type orderModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	TableID   uuid.UUID `gorm:"column:table_id;type:varchar(36)"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m orderModel) TableName() string {
	return "hueat_order"
}

func (m orderModel) toEntity() orderEntity {
	return orderEntity(m)
}

type courseModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	OrderID   uuid.UUID `gorm:"column:order_id;type:varchar(36)"`
	Position  int64     `gorm:"column:position;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m courseModel) TableName() string {
	return "hueat_course"
}

func (m courseModel) toEntity() courseEntity {
	return courseEntity(m)
}

type courseSelectionModel struct {
	ID           uuid.UUID  `gorm:"primaryKey;column:id;type:varchar(36)"`
	CourseID     uuid.UUID  `gorm:"column:course_id;type:varchar(36)"`
	MenuItemID   uuid.UUID  `gorm:"column:menu_item_id;type:varchar(36)"`
	MenuOptionID *uuid.UUID `gorm:"column:menu_option_id;type:varchar(36)"`
	Quantity     int64      `gorm:"column:quantity;type:bigint"`
	Note         *string    `gorm:"column:note;type:text"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m courseSelectionModel) TableName() string {
	return "hueat_course_selection"
}

func (m courseSelectionModel) toEntity() courseSelectionEntity {
	return courseSelectionEntity(m)
}

type orderDetailModel struct {
	Username        string    `gorm:"column:username"`
	TableName       string    `gorm:"column:table_name"`
	TableCreatedAt  time.Time `gorm:"column:table_created_at"`
	PrinterId       string    `gorm:"column:printer_id"`
	PrinterTitle    string    `gorm:"column:printer_title"`
	PrinterURL      string    `gorm:"column:printer_url"`
	CourseID        string    `gorm:"column:course_id"`
	CourseNumber    int64     `gorm:"column:course_number"`
	MenuItemTitle   string    `gorm:"column:menu_item_title"`
	MenuItemPrice   int64     `gorm:"column:menu_item_price"`
	MenuOptionTitle *string   `gorm:"column:menu_option_title"`
	MenuOptionPrice *int64    `gorm:"column:menu_option_price"`
	Quantity        int64     `gorm:"column:quantity"`
	Note            *string   `gorm:"column:note"`
}

func (m orderDetailModel) toEntity() orderDetailEntity {
	return orderDetailEntity(m)
}

type paymentDetailModel struct {
	Username       string    `gorm:"column:username"`
	TableName      string    `gorm:"column:table_name"`
	TableCreatedAt time.Time `gorm:"column:table_created_at"`
	TablePayment   string    `gorm:"column:table_payment"`
	PrinterId      string    `gorm:"column:printer_id"`
	PrinterTitle   string    `gorm:"column:printer_title"`
	PrinterURL     string    `gorm:"column:printer_url"`
	PriceTotal     int64     `gorm:"column:price_total"`
}

func (m paymentDetailModel) toEntity() paymentDetailEntity {
	return paymentDetailEntity(m)
}
