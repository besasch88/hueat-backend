package menuOption

import (
	"math"
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type menuOptionServiceInterface interface {
	listMenuOptions(ctx *gin.Context, input listMenuOptionsInputDto) ([]menuOptionEntity, int64, error)
	getMenuOptionByID(ctx *gin.Context, input getMenuOptionInputDto) (menuOptionEntity, error)
	createMenuOption(ctx *gin.Context, input createMenuOptionInputDto) (menuOptionEntity, error)
	updateMenuOption(ctx *gin.Context, input updateMenuOptionInputDto) (menuOptionEntity, error)
	deleteMenuOption(ctx *gin.Context, input deleteMenuOptionInputDto) (menuOptionEntity, error)
}

type menuOptionService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  menuOptionRepositoryInterface
}

func newMenuOptionService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository menuOptionRepositoryInterface) menuOptionService {
	return menuOptionService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s menuOptionService) listMenuOptions(ctx *gin.Context, input listMenuOptionsInputDto) ([]menuOptionEntity, int64, error) {
	menuItemID := uuid.MustParse(input.MenuItemId)
	if exists, err := s.repository.checkMenuItemExists(s.storage, menuItemID); err != nil {
		return []menuOptionEntity{}, 0, hueat_err.ErrGeneric
	} else if !exists {
		return []menuOptionEntity{}, 0, errMenuItemNotFound
	}
	items, totalCount, err := s.repository.listMenuOptions(s.storage, menuItemID, false)
	if err != nil || items == nil {
		return []menuOptionEntity{}, 0, hueat_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s menuOptionService) getMenuOptionByID(ctx *gin.Context, input getMenuOptionInputDto) (menuOptionEntity, error) {
	menuOptionID := uuid.MustParse(input.ID)
	item, err := s.repository.getMenuOptionByID(s.storage, menuOptionID, false)
	if err != nil {
		return menuOptionEntity{}, hueat_err.ErrGeneric
	}
	if hueat_utils.IsEmpty(item) {
		return menuOptionEntity{}, errMenuOptionNotFound
	}
	return item, nil
}

func (s menuOptionService) createMenuOption(ctx *gin.Context, input createMenuOptionInputDto) (menuOptionEntity, error) {
	menuItemId := uuid.MustParse(input.MenuItemId)
	menuItemID := uuid.MustParse(input.MenuItemId)
	if exists, err := s.repository.checkMenuItemExists(s.storage, menuItemID); err != nil {
		return menuOptionEntity{}, hueat_err.ErrGeneric
	} else if !exists {
		return menuOptionEntity{}, errMenuItemNotFound
	}
	now := time.Now()
	maxValue := int64(math.MaxInt64)
	newMenuOption := menuOptionEntity{
		ID:         uuid.New(),
		MenuItemID: menuItemId,
		Position:   maxValue,
		Title:      input.Title,
		Active:     hueat_utils.BoolPtr(false),
		Inside:     hueat_utils.BoolPtr(true),
		Outside:    hueat_utils.BoolPtr(true),
		Price:      input.Price,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuOptionEntity
		if item, err := s.repository.getMenuOptionByTitle(tx, input.Title, false); err != nil {
			return hueat_err.ErrGeneric
		} else if !hueat_utils.IsEmpty(item) {
			return errMenuOptionSameTitleAlreadyExists
		} else if _, err = s.repository.saveMenuOption(tx, newMenuOption, hueat_db.Create); err != nil {
			return hueat_err.ErrGeneric
		} else if updatedEntities, err = s.repository.recalculateMenuOptionsPosition(tx, menuItemId); err != nil {
			return hueat_err.ErrGeneric
		} else if newMenuOption, err = s.repository.getMenuOptionByID(tx, newMenuOption.ID, false); err != nil {
			return hueat_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuOptionV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.MenuOptionCreatedEvent,
				EventEntity: &hueat_pubsub.MenuOptionEventEntity{
					ID:        newMenuOption.ID,
					Title:     newMenuOption.Title,
					Position:  newMenuOption.Position,
					Active:    newMenuOption.Active,
					Inside:    newMenuOption.Inside,
					Outside:   newMenuOption.Outside,
					Price:     newMenuOption.Price,
					CreatedAt: newMenuOption.CreatedAt,
					UpdatedAt: newMenuOption.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(menuOptionEntity{}, newMenuOption),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == newMenuOption.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuOptionV1, hueat_pubsub.PubSubMessage{
				Message: hueat_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: hueat_pubsub.MenuOptionUpdatedEvent,
					EventEntity: &hueat_pubsub.MenuOptionEventEntity{
						ID:         updatedEntity.ID,
						MenuItemID: updatedEntity.MenuItemID,
						Title:      updatedEntity.Title,
						Position:   updatedEntity.Position,
						Active:     updatedEntity.Active,
						Inside:     updatedEntity.Inside,
						Outside:    updatedEntity.Outside,
						Price:      updatedEntity.Price,
						CreatedAt:  updatedEntity.CreatedAt,
						UpdatedAt:  updatedEntity.UpdatedAt,
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
		return menuOptionEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newMenuOption, nil
}

func (s menuOptionService) updateMenuOption(ctx *gin.Context, input updateMenuOptionInputDto) (menuOptionEntity, error) {
	now := time.Now()
	var updatedMenuOption menuOptionEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuOptionEntity
		menuOptionId := uuid.MustParse(input.ID)
		currentMenuOption, err := s.repository.getMenuOptionByID(tx, menuOptionId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentMenuOption) {
			return errMenuOptionNotFound
		}
		updatedMenuOption = currentMenuOption
		updatedMenuOption.UpdatedAt = now
		// If the input contains a new title, check for collision
		if input.Title != nil {
			menuOptionSameTitle, err := s.repository.getMenuOptionByTitle(tx, *input.Title, false)
			if err != nil {
				return hueat_err.ErrGeneric
			}
			if !hueat_utils.IsEmpty(menuOptionSameTitle) && menuOptionSameTitle.ID.String() != input.ID {
				return errMenuOptionSameTitleAlreadyExists
			}
			updatedMenuOption.Title = *input.Title
		}
		if input.Active != nil {
			updatedMenuOption.Active = input.Active
		}
		if input.Inside != nil {
			updatedMenuOption.Inside = input.Inside
		}
		if input.Outside != nil {
			updatedMenuOption.Outside = input.Outside
		}
		if input.Price != nil {
			updatedMenuOption.Price = *input.Price
		}
		if input.Position != nil {
			// If the step is moving in a lower position (e.g. from 10 to 3),
			// we need to move it one step more, so that, the algorith to re-sort all steps correctly
			if updatedMenuOption.Position < *input.Position {
				*input.Position++
			}
			updatedMenuOption.Position = *input.Position
		}
		if _, err = s.repository.saveMenuOption(tx, updatedMenuOption, hueat_db.Update); err != nil {
			return hueat_err.ErrGeneric
		}
		if updatedEntities, err = s.repository.recalculateMenuOptionsPosition(tx, updatedMenuOption.MenuItemID); err != nil {
			return hueat_err.ErrGeneric
		}
		if updatedMenuOption, err = s.repository.getMenuOptionByID(tx, updatedMenuOption.ID, false); err != nil {
			return hueat_err.ErrGeneric
		}

		// Send an event of menuOption updated
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuOptionV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.MenuOptionUpdatedEvent,
				EventEntity: &hueat_pubsub.MenuOptionEventEntity{
					ID:         updatedMenuOption.ID,
					MenuItemID: updatedMenuOption.MenuItemID,
					Title:      updatedMenuOption.Title,
					Position:   updatedMenuOption.Position,
					Active:     updatedMenuOption.Active,
					Inside:     updatedMenuOption.Inside,
					Outside:    updatedMenuOption.Outside,
					Price:      updatedMenuOption.Price,
					CreatedAt:  updatedMenuOption.CreatedAt,
					UpdatedAt:  updatedMenuOption.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentMenuOption, updatedMenuOption),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if updatedEntity.ID == updatedMenuOption.ID {
				continue
			}
			if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuOptionV1, hueat_pubsub.PubSubMessage{
				Message: hueat_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: hueat_pubsub.MenuOptionUpdatedEvent,
					EventEntity: &hueat_pubsub.MenuOptionEventEntity{
						ID:         updatedEntity.ID,
						MenuItemID: updatedEntity.MenuItemID,
						Title:      updatedEntity.Title,
						Position:   updatedEntity.Position,
						Active:     updatedEntity.Active,
						Inside:     updatedEntity.Inside,
						Outside:    updatedEntity.Outside,
						Price:      updatedEntity.Price,
						CreatedAt:  updatedEntity.CreatedAt,
						UpdatedAt:  updatedEntity.UpdatedAt,
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
		return menuOptionEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedMenuOption, nil
}

func (s menuOptionService) deleteMenuOption(ctx *gin.Context, input deleteMenuOptionInputDto) (menuOptionEntity, error) {
	var currentMenuOption menuOptionEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		var updatedEntities []menuOptionEntity
		// Check if the Menu Item exists
		menuOptionId := uuid.MustParse(input.ID)
		currentMenuOption, err := s.repository.getMenuOptionByID(tx, menuOptionId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentMenuOption) {
			return errMenuOptionNotFound
		}
		s.repository.deleteMenuOption(tx, currentMenuOption)
		if updatedEntities, err = s.repository.recalculateMenuOptionsPosition(tx, currentMenuOption.MenuItemID); err != nil {
			return hueat_err.ErrGeneric
		}
		// Send an event of menuOption deleted
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuOptionV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.MenuOptionDeletedEvent,
				EventEntity: &hueat_pubsub.MenuOptionEventEntity{
					ID:         currentMenuOption.ID,
					MenuItemID: currentMenuOption.MenuItemID,
					Title:      currentMenuOption.Title,
					Position:   currentMenuOption.Position,
					Active:     currentMenuOption.Active,
					Inside:     currentMenuOption.Inside,
					Outside:    currentMenuOption.Outside,
					Price:      currentMenuOption.Price,
					CreatedAt:  currentMenuOption.CreatedAt,
					UpdatedAt:  currentMenuOption.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentMenuOption, menuOptionEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		// For the list of updated entities in Position, send events
		for _, updatedEntity := range updatedEntities {
			if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicMenuItemV1, hueat_pubsub.PubSubMessage{
				Message: hueat_pubsub.PubSubEvent{
					EventID:   uuid.New(),
					EventTime: time.Now(),
					EventType: hueat_pubsub.MenuOptionUpdatedEvent,
					EventEntity: &hueat_pubsub.MenuOptionEventEntity{
						ID:         updatedEntity.ID,
						MenuItemID: updatedEntity.MenuItemID,
						Title:      updatedEntity.Title,
						Position:   updatedEntity.Position,
						Active:     updatedEntity.Active,
						Inside:     updatedEntity.Inside,
						Outside:    updatedEntity.Outside,
						CreatedAt:  updatedEntity.CreatedAt,
						UpdatedAt:  updatedEntity.UpdatedAt,
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
		return menuOptionEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentMenuOption, nil
}
