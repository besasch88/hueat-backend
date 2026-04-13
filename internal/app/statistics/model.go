package statistics

import (
	"time"

	"github.com/google/uuid"
)

type tableModel struct {
	ID uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
}

func (m tableModel) TableName() string {
	return "hueat_table"
}

type averageTableDurationModel struct {
	Duration time.Duration `gorm:"column:avg_duration"`
}

type paymentTakingsModel struct {
	PaymentType string `gorm:"column:payment_method"`
	Takings     int64  `gorm:"column:takings"`
}

func (m paymentTakingsModel) toEntity() paymentTakingsEntity {
	return paymentTakingsEntity(m)
}

type menuItemStatModel struct {
	Title    string `gorm:"column:title"`
	Quantity int64  `gorm:"column:quantity"`
	Takings  int64  `gorm:"column:takings"`
}

func (m menuItemStatModel) toEntity() menuItemStatEntity {
	return menuItemStatEntity(m)
}
