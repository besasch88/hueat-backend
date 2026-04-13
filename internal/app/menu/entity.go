package menu

import (
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
)

type menuCategoryEntity hueat_pubsub.MenuCategoryEventEntity
type menuItemEntity hueat_pubsub.MenuItemEventEntity
type menuOptionEntity hueat_pubsub.MenuOptionEventEntity

type menuOption struct {
	menuOptionEntity
}
type menuItem struct {
	menuItemEntity
	Options []menuOption `json:"options"`
}
type menuCategory struct {
	menuCategoryEntity
	Items []menuItem `json:"items"`
}

type menu struct {
	Categories []menuCategory `json:"categories"`
}
