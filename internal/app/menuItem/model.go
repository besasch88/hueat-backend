package menuItem

import (
	"time"

	"github.com/google/uuid"
)

type menuCategoryModel struct {
	ID uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
}

func (m menuCategoryModel) TableName() string {
	return "hueat_menu_category"
}

type menuItemModel struct {
	ID             uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	MenuCategoryID uuid.UUID `gorm:"column:menu_category_id;type:varchar(36)"`
	Title          string    `gorm:"column:title;type:varchar(255)"`
	Position       int64     `gorm:"column:position;type:bigint"`
	Active         *bool     `gorm:"column:active;type:boolean"`
	Inside         *bool     `gorm:"column:inside;type:boolean"`
	Outside        *bool     `gorm:"column:outside;type:boolean"`
	Price          int64     `gorm:"column:price;type:bigint"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime:false"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false"`
}

func (m menuItemModel) TableName() string {
	return "hueat_menu_item"
}

func (m menuItemModel) toEntity() menuItemEntity {
	return menuItemEntity(m)
}
