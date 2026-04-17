package menuItem

import (
	"github.com/google/uuid"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
)

type menuItemEntity hueat_pubsub.MenuItemEventEntity

type tableEntity struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
	Close  *bool     `json:"close"`
	Inside *bool     `json:"inside"`
}
