package menuItem

import (
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"gorm.io/gorm"
)

type menuItemServiceInterface interface {
	listMenuItems(ctx *gin.Context, input listMenuItemsInputDto) ([]menuItemEntity, int64, error)
	getMenuItemByID(ctx *gin.Context, input getMenuItemInputDto) (menuItemEntity, error)
	createMenuItem(ctx *gin.Context, input createMenuItemInputDto) (menuItemEntity, error)
	updateMenuItem(ctx *gin.Context, input updateMenuItemInputDto) (menuItemEntity, error)
	deleteMenuItem(ctx *gin.Context, input deleteMenuItemInputDto) (menuItemEntity, error)
}

type menuItemService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  menuItemRepositoryInterface
}

func newMenuItemService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository menuItemRepositoryInterface) menuItemService {
	return menuItemService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s menuItemService) listMenuItems(ctx *gin.Context, input listMenuItemsInputDto) ([]menuItemEntity, int64, error) {
	menuCategoryID := uuid.MustParse(input.MenuCategoryId)
	if exists, err := s.repository.checkMenuCategoryExists(s.storage, menuCategoryID); err != nil {
		return []menuItemEntity{}, 0, hueat_err.ErrGeneric
	} else if !exists {
		return []menuItemEntity{}, 0, errMenuCategoryNotFound
	}
	items, totalCount, err := s.repository.listMenuItems(s.storage, menuCategoryID, false)
	if err != nil || items == nil {
		return []menuItemEntity{}, 0, hueat_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s menuItemService) getMenuItemByID(ctx *gin.Context, input getMenuItemInputDto) (menuItemEntity, error) {
	menuItemID := uuid.MustParse(input.ID)
	item, err := s.repository.getMenuItemByID(s.storage, menuItemID, false)
	if err != nil {
		return menuItemEntity{}, hueat_err.ErrGeneric
	}
	if hueat_utils.IsEmpty(item) {
		return menuItemEntity{}, errMenuItemNotFound
	}
	return item, nil
}

func (s menuItemService) createMenuItem(ctx *gin.Context, input createMenuItemInputDto) (menuItemEntity, error) {
	menuCategoryID := uuid.MustParse(input.MenuCategoryId)
	if input.PrinterInsideID != nil {
		printerId := uuid.MustParse(*input.PrinterInsideID)
		if exists, err := s.repository.checkPrinterExists(s.storage, printerId); err != nil {
			return menuItemEntity{}, hueat_err.ErrGeneric
		} else if !exists {
			return menuItemEntity{}, errPrinterNotFound
		}
	}
	if input.PrinterOutsideID != nil {
		printerId := uuid.MustParse(*input.PrinterOutsideID)
		if exists, err := s.repository.checkPrinterExists(s.storage, printerId); err != nil {
			return menuItemEntity{}, hueat_err.ErrGeneric
		} else if !exists {
			return menuItemEntity{}, errPrinterNotFound
		}
	}
	if exists, err := s.repository.checkMenuCategoryExists(s.storage, menuCategoryID); err != nil {
		return menuItemEntity{}, hueat_err.ErrGeneric
	} else if !exists {
		return menuItemEntity{}, errMenuCategoryNotFound
	}
	now := time.Now()
	maxValue := int64(math.MaxInt64)
	newMenuItem := menuItemEntity{
		ID:               uuid.New(),
		MenuCategoryID:   menuCategoryID,
		Position:         maxValue,
		Title:            input.Title,
		TitleDisplay:     input.TitleDisplay,
		Active:           hueat_utils.BoolPtr(false),
		Inside:           hueat_utils.BoolPtr(true),
		Outside:          hueat_utils.BoolPtr(true),
		Price:            input.Price,
		PrinterInsideID:  hueat_utils.GetOptionalUUIDFromString(input.PrinterInsideID),
		PrinterOutsideID: hueat_utils.GetOptionalUUIDFromString(input.PrinterOutsideID),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuItemEntity
		if item, err := s.repository.getMenuItemByTitle(tx, input.Title, false); err != nil {
			return hueat_err.ErrGeneric
		} else if !hueat_utils.IsEmpty(item) {
			return errMenuItemSameTitleAlreadyExists
		} else if _, err = s.repository.saveMenuItem(tx, newMenuItem, hueat_db.Create); err != nil {
			return hueat_err.ErrGeneric
		} else if updatedEntities, err = s.repository.recalculateMenuItemsPosition(tx, menuCategoryID); err != nil {
			return hueat_err.ErrGeneric
		} else if newMenuItem, err = s.repository.getMenuItemByID(tx, newMenuItem.ID, false); err != nil {
			return hueat_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuItemV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.MenuItemCreatedEvent,
				EventEntity: &hueat_pubsub.MenuItemEventEntity{
					ID:               newMenuItem.ID,
					Title:            newMenuItem.Title,
					TitleDisplay:     newMenuItem.TitleDisplay,
					Position:         newMenuItem.Position,
					Active:           newMenuItem.Active,
					Inside:           newMenuItem.Inside,
					Outside:          newMenuItem.Outside,
					Price:            newMenuItem.Price,
					PrinterInsideID:  newMenuItem.PrinterInsideID,
					PrinterOutsideID: newMenuItem.PrinterOutsideID,
					CreatedAt:        newMenuItem.CreatedAt,
					UpdatedAt:        newMenuItem.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(menuItemEntity{}, newMenuItem),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == newMenuItem.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuItemV1, hueat_pubsub.PubSubMessage{
				Message: hueat_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: hueat_pubsub.MenuItemUpdatedEvent,
					EventEntity: &hueat_pubsub.MenuItemEventEntity{
						ID:               updatedEntity.ID,
						MenuCategoryID:   updatedEntity.MenuCategoryID,
						Title:            updatedEntity.Title,
						TitleDisplay:     updatedEntity.TitleDisplay,
						Position:         updatedEntity.Position,
						Active:           updatedEntity.Active,
						Inside:           updatedEntity.Inside,
						Outside:          updatedEntity.Outside,
						Price:            updatedEntity.Price,
						PrinterInsideID:  updatedEntity.PrinterInsideID,
						PrinterOutsideID: updatedEntity.PrinterOutsideID,
						CreatedAt:        updatedEntity.CreatedAt,
						UpdatedAt:        updatedEntity.UpdatedAt,
					},
					EventChangedFields: []string{"Position", "UpdatedAt"},
				},
			}); err != nil {
				return err
			} else {
				eventsToPublish = append(eventsToPublish, event)
			}
		}
		return nil
	})
	if errTransaction != nil {
		return menuItemEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newMenuItem, nil
}

func (s menuItemService) updateMenuItem(ctx *gin.Context, input updateMenuItemInputDto) (menuItemEntity, error) {
	if input.PrinterInsideID != nil {
		printerId := uuid.MustParse(*input.PrinterInsideID)
		if exists, err := s.repository.checkPrinterExists(s.storage, printerId); err != nil {
			return menuItemEntity{}, hueat_err.ErrGeneric
		} else if !exists {
			return menuItemEntity{}, errPrinterNotFound
		}
	}
	if input.PrinterOutsideID != nil {
		printerId := uuid.MustParse(*input.PrinterOutsideID)
		if exists, err := s.repository.checkPrinterExists(s.storage, printerId); err != nil {
			return menuItemEntity{}, hueat_err.ErrGeneric
		} else if !exists {
			return menuItemEntity{}, errPrinterNotFound
		}
	}
	now := time.Now()
	var updatedMenuItem menuItemEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuItemEntity
		menuItemId := uuid.MustParse(input.ID)
		currentMenuItem, err := s.repository.getMenuItemByID(tx, menuItemId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentMenuItem) {
			return errMenuItemNotFound
		}
		updatedMenuItem = currentMenuItem
		updatedMenuItem.UpdatedAt = now
		// If the input contains a new title, check for collision
		if input.Title != nil {
			menuItemSameTitle, err := s.repository.getMenuItemByTitle(tx, *input.Title, false)
			if err != nil {
				return hueat_err.ErrGeneric
			}
			if !hueat_utils.IsEmpty(menuItemSameTitle) && menuItemSameTitle.ID.String() != input.ID {
				return errMenuItemSameTitleAlreadyExists
			}
			updatedMenuItem.Title = *input.Title
		}
		if input.TitleDisplay != nil {
			updatedMenuItem.TitleDisplay = *input.TitleDisplay
		}
		if input.Active != nil {
			updatedMenuItem.Active = input.Active
		}
		if input.Inside != nil {
			updatedMenuItem.Inside = input.Inside
		}
		if input.Outside != nil {
			updatedMenuItem.Outside = input.Outside
		}
		if input.Price != nil {
			updatedMenuItem.Price = *input.Price
		}
		if input.PrinterInsideID != nil {
			updatedMenuItem.PrinterInsideID = hueat_utils.GetOptionalUUIDFromString(input.PrinterInsideID)
		}
		if input.PrinterOutsideID != nil {
			updatedMenuItem.PrinterOutsideID = hueat_utils.GetOptionalUUIDFromString(input.PrinterOutsideID)
		}
		if input.Position != nil {
			// If the step is moving in a lower position (e.g. from 10 to 3),
			// we need to move it one step more, so that, the algorith to re-sort all steps correctly
			if updatedMenuItem.Position < *input.Position {
				*input.Position++
			}
			updatedMenuItem.Position = *input.Position
		}
		if _, err = s.repository.saveMenuItem(tx, updatedMenuItem, hueat_db.Update); err != nil {
			return hueat_err.ErrGeneric
		}
		if updatedEntities, err = s.repository.recalculateMenuItemsPosition(tx, updatedMenuItem.MenuCategoryID); err != nil {
			return hueat_err.ErrGeneric
		}
		if updatedMenuItem, err = s.repository.getMenuItemByID(tx, updatedMenuItem.ID, false); err != nil {
			return hueat_err.ErrGeneric
		}

		// Send an event of menuItem updated
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuItemV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.MenuItemUpdatedEvent,
				EventEntity: &hueat_pubsub.MenuItemEventEntity{
					ID:               updatedMenuItem.ID,
					MenuCategoryID:   updatedMenuItem.MenuCategoryID,
					Title:            updatedMenuItem.Title,
					TitleDisplay:     updatedMenuItem.TitleDisplay,
					Position:         updatedMenuItem.Position,
					Active:           updatedMenuItem.Active,
					Inside:           updatedMenuItem.Inside,
					Outside:          updatedMenuItem.Outside,
					Price:            updatedMenuItem.Price,
					PrinterInsideID:  updatedMenuItem.PrinterInsideID,
					PrinterOutsideID: updatedMenuItem.PrinterOutsideID,
					CreatedAt:        updatedMenuItem.CreatedAt,
					UpdatedAt:        updatedMenuItem.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentMenuItem, updatedMenuItem),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == updatedMenuItem.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuItemV1, hueat_pubsub.PubSubMessage{
				Message: hueat_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: hueat_pubsub.MenuItemUpdatedEvent,
					EventEntity: &hueat_pubsub.MenuItemEventEntity{
						ID:               updatedEntity.ID,
						MenuCategoryID:   updatedEntity.MenuCategoryID,
						Title:            updatedEntity.Title,
						TitleDisplay:     updatedEntity.TitleDisplay,
						Position:         updatedEntity.Position,
						Active:           updatedEntity.Active,
						Inside:           updatedEntity.Inside,
						Outside:          updatedEntity.Outside,
						Price:            updatedEntity.Price,
						PrinterInsideID:  updatedEntity.PrinterInsideID,
						PrinterOutsideID: updatedEntity.PrinterOutsideID,
						CreatedAt:        updatedEntity.CreatedAt,
						UpdatedAt:        updatedEntity.UpdatedAt,
					},
					EventChangedFields: []string{"Position", "UpdatedAt"},
				},
			}); err != nil {
				return err
			} else {
				eventsToPublish = append(eventsToPublish, event)
			}
		}
		return nil
	})
	if errTransaction != nil {
		return menuItemEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedMenuItem, nil
}

func (s menuItemService) deleteMenuItem(ctx *gin.Context, input deleteMenuItemInputDto) (menuItemEntity, error) {
	var currentMenuItem menuItemEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuItemEntity
		// Check if the Menu Item exists
		menuItemId := uuid.MustParse(input.ID)
		currentMenuItem, err := s.repository.getMenuItemByID(tx, menuItemId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentMenuItem) {
			return errMenuItemNotFound
		}
		s.repository.deleteMenuItem(tx, currentMenuItem)
		if updatedEntities, err = s.repository.recalculateMenuItemsPosition(tx, currentMenuItem.MenuCategoryID); err != nil {
			return hueat_err.ErrGeneric
		}
		// Send an event of menuItem deleted
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuItemV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.MenuItemDeletedEvent,
				EventEntity: &hueat_pubsub.MenuItemEventEntity{
					ID:               currentMenuItem.ID,
					MenuCategoryID:   currentMenuItem.MenuCategoryID,
					Title:            currentMenuItem.Title,
					TitleDisplay:     currentMenuItem.TitleDisplay,
					Position:         currentMenuItem.Position,
					Active:           currentMenuItem.Active,
					Inside:           currentMenuItem.Inside,
					Outside:          currentMenuItem.Outside,
					Price:            currentMenuItem.Price,
					PrinterInsideID:  currentMenuItem.PrinterInsideID,
					PrinterOutsideID: currentMenuItem.PrinterOutsideID,
					CreatedAt:        currentMenuItem.CreatedAt,
					UpdatedAt:        currentMenuItem.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentMenuItem, menuItemEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuCategoryV1, hueat_pubsub.PubSubMessage{
				Message: hueat_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: hueat_pubsub.MenuItemUpdatedEvent,
					EventEntity: &hueat_pubsub.MenuItemEventEntity{
						ID:               updatedEntity.ID,
						MenuCategoryID:   updatedEntity.MenuCategoryID,
						Title:            updatedEntity.Title,
						TitleDisplay:     updatedEntity.TitleDisplay,
						Position:         updatedEntity.Position,
						Active:           updatedEntity.Active,
						Inside:           updatedEntity.Inside,
						Outside:          updatedEntity.Outside,
						Price:            updatedEntity.Price,
						PrinterInsideID:  updatedEntity.PrinterInsideID,
						PrinterOutsideID: updatedEntity.PrinterOutsideID,
						CreatedAt:        updatedEntity.CreatedAt,
						UpdatedAt:        updatedEntity.UpdatedAt,
					},
					EventChangedFields: []string{"Position", "UpdatedAt"},
				},
			}); err != nil {
				return err
			} else {
				eventsToPublish = append(eventsToPublish, event)
			}
		}
		return nil
	})
	if errTransaction != nil {
		return menuItemEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentMenuItem, nil
}
