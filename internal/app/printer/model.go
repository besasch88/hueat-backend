package printer

import (
	"time"

	"github.com/google/uuid"
)

type printerModel struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	Title     string    `gorm:"column:title;type:varchar(255)"`
	Url       string    `gorm:"column:url;type:varchar(255)"`
	Active    *bool     `gorm:"column:active;type:boolean"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m printerModel) TableName() string {
	return "hueat_printer"
}

func (m printerModel) toEntity() printerEntity {
	return printerEntity(m)
}
