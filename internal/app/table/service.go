package table

import (
	"slices"
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tableServiceInterface interface {
	listTables(ctx *gin.Context, input listTablesInputDto) ([]tableEntity, int64, error)
	getTableByID(ctx *gin.Context, input getTableInputDto) (tableEntity, error)
	createTable(ctx *gin.Context, input createTableInputDto) (tableEntity, error)
	updateTable(ctx *gin.Context, input updateTableInputDto) (tableEntity, error)
	deleteTable(ctx *gin.Context, input deleteTableInputDto) (tableEntity, error)
}

type tableService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  tableRepositoryInterface
}

func newTableService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository tableRepositoryInterface) tableService {
	return tableService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s tableService) listTables(ctx *gin.Context, input listTablesInputDto) ([]tableEntity, int64, error) {
	requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
	userId := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, hueat_auth.READ_OTHER_TABLES) {
		userId = nil
	}
	items, totalCount, err := s.repository.listTables(s.storage, userId, input.Target == "inside", input.IncludeClosed, false)
	if err != nil || items == nil {
		return []tableEntity{}, 0, hueat_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s tableService) getTableByID(ctx *gin.Context, input getTableInputDto) (tableEntity, error) {
	requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
	userId := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, hueat_auth.READ_OTHER_TABLES) {
		userId = nil
	}
	tableID := uuid.MustParse(input.ID)
	item, err := s.repository.getTableByID(s.storage, tableID, userId, false)
	if err != nil {
		return tableEntity{}, hueat_err.ErrGeneric
	}
	if hueat_utils.IsEmpty(item) {
		return tableEntity{}, errTableNotFound
	}
	return item, nil
}

func (s tableService) createTable(ctx *gin.Context, input createTableInputDto) (tableEntity, error) {
	requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
	userId := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
	now := time.Now()
	newTable := tableEntity{
		ID:            uuid.New(),
		UserID:        *userId,
		Name:          input.Name,
		Inside:        input.Inside,
		Close:         hueat_utils.BoolPtr(false),
		PaymentMethod: nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		if item, err := s.repository.getOpenTableByName(tx, input.Name, false); err != nil {
			return hueat_err.ErrGeneric
		} else if !hueat_utils.IsEmpty(item) {
			return errTableSameNameAlreadyExists
		} else if _, err = s.repository.saveTable(tx, newTable, hueat_db.Create); err != nil {
			return hueat_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicTableV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.TableCreatedEvent,
				EventEntity: &hueat_pubsub.TableEventEntity{
					ID:            newTable.ID,
					UserID:        newTable.UserID,
					Name:          newTable.Name,
					Inside:        newTable.Inside,
					Close:         newTable.Close,
					PaymentMethod: newTable.PaymentMethod,
					CreatedAt:     newTable.CreatedAt,
					UpdatedAt:     newTable.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(tableEntity{}, newTable),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return tableEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newTable, nil
}

func (s tableService) updateTable(ctx *gin.Context, input updateTableInputDto) (tableEntity, error) {
	requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
	userId := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, hueat_auth.WRITE_OTHER_TABLES) {
		userId = nil
	}
	now := time.Now()
	var updatedTable tableEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		tableId := uuid.MustParse(input.ID)
		currentTable, err := s.repository.getTableByID(tx, tableId, userId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentTable) {
			return errTableNotFound
		}
		updatedTable = currentTable
		updatedTable.UpdatedAt = now
		// If the input contains a new name, check for collision
		if input.Name != nil {
			tableSameName, err := s.repository.getOpenTableByName(tx, *input.Name, false)
			if err != nil {
				return hueat_err.ErrGeneric
			}
			if !hueat_utils.IsEmpty(tableSameName) && tableSameName.ID.String() != input.ID {
				return errTableSameNameAlreadyExists
			}
			updatedTable.Name = *input.Name
		}
		if input.Inside != nil {
			updatedTable.Inside = input.Inside
		}
		if input.Close != nil {
			updatedTable.Close = input.Close
		}
		if input.PaymentMethod != nil {
			updatedTable.PaymentMethod = input.PaymentMethod
		}
		if _, err = s.repository.saveTable(tx, updatedTable, hueat_db.Update); err != nil {
			return hueat_err.ErrGeneric
		}
		if updatedTable, err = s.repository.getTableByID(tx, updatedTable.ID, userId, false); err != nil {
			return hueat_err.ErrGeneric
		}

		// Send an event of table updated
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicTableV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.TableUpdatedEvent,
				EventEntity: &hueat_pubsub.TableEventEntity{
					ID:            updatedTable.ID,
					UserID:        updatedTable.UserID,
					Name:          updatedTable.Name,
					Inside:        updatedTable.Inside,
					Close:         updatedTable.Close,
					PaymentMethod: updatedTable.PaymentMethod,
					CreatedAt:     updatedTable.CreatedAt,
					UpdatedAt:     updatedTable.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentTable, updatedTable),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return tableEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedTable, nil
}

func (s tableService) deleteTable(ctx *gin.Context, input deleteTableInputDto) (tableEntity, error) {
	requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
	userId := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, hueat_auth.WRITE_OTHER_TABLES) {
		userId = nil
	}
	var currentTable tableEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		// Check if the Menu Item exists
		tableId := uuid.MustParse(input.ID)
		currentTable, err := s.repository.getTableByID(tx, tableId, userId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentTable) {
			return errTableNotFound
		}
		s.repository.deleteTable(tx, currentTable)
		// Send an event of table deleted
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicTableV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.TableDeletedEvent,
				EventEntity: &hueat_pubsub.TableEventEntity{
					ID:            currentTable.ID,
					UserID:        currentTable.UserID,
					Name:          currentTable.Name,
					Inside:        currentTable.Inside,
					Close:         currentTable.Close,
					PaymentMethod: currentTable.PaymentMethod,
					CreatedAt:     currentTable.CreatedAt,
					UpdatedAt:     currentTable.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentTable, tableEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return tableEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentTable, nil
}
