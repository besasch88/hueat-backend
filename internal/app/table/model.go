package table

import (
	"time"

	"github.com/google/uuid"
)

type tableModel struct {
	ID            uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	UserID        uuid.UUID `gorm:"column:user_id;type:varchar(36)"`
	Name          string    `gorm:"column:name;type:varchar(255)"`
	Close         *bool     `gorm:"column:close;type:boolean"`
	Inside        *bool     `gorm:"column:inside;type:boolean"`
	PaymentMethod *string   `gorm:"column:payment_method;type:varchar(255)"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m tableModel) TableName() string {
	return "hueat_table"
}

func (m tableModel) toEntity() tableEntity {
	return tableEntity(m)
}
