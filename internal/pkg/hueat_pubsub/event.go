package hueat_pubsub

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

/*
PubSubEventType represents an event type that can be published or consumed within the pub-sub system.
Generally, the PubSubEventType is related to an event entity and the possible actions performed on it.
It is preferable to use the past participle to indicate that the event was generated as a result
of an application state change.
*/
type PubSubEventType string

/*
List of avaiable events can be published and consumed within the pub-sub system.
*/
const (
	PrinterCreatedEvent      PubSubEventType = "printer.created"
	PrinterUpdatedEvent      PubSubEventType = "printer.updated"
	PrinterDeletedEvent      PubSubEventType = "printer.deleted"
	MenuCategoryCreatedEvent PubSubEventType = "menu-category.created"
	MenuCategoryUpdatedEvent PubSubEventType = "menu-category.updated"
	MenuCategoryDeletedEvent PubSubEventType = "menu-category.deleted"
	MenuItemCreatedEvent     PubSubEventType = "menu-item.created"
	MenuItemUpdatedEvent     PubSubEventType = "menu-item.updated"
	MenuItemDeletedEvent     PubSubEventType = "menu-item.deleted"
	MenuOptionCreatedEvent   PubSubEventType = "menu-option.created"
	MenuOptionUpdatedEvent   PubSubEventType = "menu-option.updated"
	MenuOptionDeletedEvent   PubSubEventType = "menu-option.deleted"
	TableCreatedEvent        PubSubEventType = "table.created"
	TableUpdatedEvent        PubSubEventType = "table.updated"
	TableDeletedEvent        PubSubEventType = "table.deleted"
	OrderCreatedEvent        PubSubEventType = "order.created"
	CourseCreatedEvent       PubSubEventType = "course.created"
)

/*
Map each event type to a function that returns a pointer to the right struct.
It is useful for unmarshal stored events and replay
*/
var eventEntityFactories = map[PubSubEventType]func() any{
	PrinterCreatedEvent:      func() interface{} { return &PrinterEventEntity{} },
	PrinterUpdatedEvent:      func() interface{} { return &PrinterEventEntity{} },
	PrinterDeletedEvent:      func() interface{} { return &PrinterEventEntity{} },
	MenuCategoryCreatedEvent: func() interface{} { return &MenuCategoryEventEntity{} },
	MenuCategoryUpdatedEvent: func() interface{} { return &MenuCategoryEventEntity{} },
	MenuCategoryDeletedEvent: func() interface{} { return &MenuCategoryEventEntity{} },
	MenuItemCreatedEvent:     func() interface{} { return &MenuItemEventEntity{} },
	MenuItemUpdatedEvent:     func() interface{} { return &MenuItemEventEntity{} },
	MenuItemDeletedEvent:     func() interface{} { return &MenuItemEventEntity{} },
	MenuOptionCreatedEvent:   func() interface{} { return &MenuOptionEventEntity{} },
	MenuOptionUpdatedEvent:   func() interface{} { return &MenuOptionEventEntity{} },
	MenuOptionDeletedEvent:   func() interface{} { return &MenuOptionEventEntity{} },
	TableCreatedEvent:        func() interface{} { return &TableEventEntity{} },
	TableUpdatedEvent:        func() interface{} { return &TableEventEntity{} },
	TableDeletedEvent:        func() interface{} { return &TableEventEntity{} },
	OrderCreatedEvent:        func() interface{} { return &OrderEventEntity{} },
	CourseCreatedEvent:       func() interface{} { return &CourseEventEntity{} },
}

/*
PubSubEvent represents a generic struct for events. All the events must be structured in this way,
ensuring the payload of the event itself is stored inside the EventEntity.
*/
type PubSubEvent struct {
	EventID            uuid.UUID       `json:"eventId"`
	EventTime          time.Time       `json:"eventTime"`
	EventType          PubSubEventType `json:"eventType"`
	EventEntity        interface{}     `json:"eventEntity"`
	EventChangedFields []string        `json:"eventChangedFields"`
	EventState         *sync.WaitGroup `json:"-"`
}
