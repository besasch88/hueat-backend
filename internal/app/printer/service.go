package printer

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type printerServiceInterface interface {
	listPrinters(ctx *gin.Context) ([]printerEntity, int64, error)
	getPrinterByID(ctx *gin.Context, input getPrinterInputDto) (printerEntity, error)
	createPrinter(ctx *gin.Context, input createPrinterInputDto) (printerEntity, error)
	updatePrinter(ctx *gin.Context, input updatePrinterInputDto) (printerEntity, error)
	deletePrinter(ctx *gin.Context, input deletePrinterInputDto) (printerEntity, error)
}

type printerService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  printerRepositoryInterface
}

func newPrinterService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository printerRepositoryInterface) printerService {
	return printerService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s printerService) listPrinters(ctx *gin.Context) ([]printerEntity, int64, error) {
	items, totalCount, err := s.repository.listPrinters(s.storage, false)
	if err != nil || items == nil {
		return []printerEntity{}, 0, hueat_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s printerService) getPrinterByID(ctx *gin.Context, input getPrinterInputDto) (printerEntity, error) {
	printerID := uuid.MustParse(input.ID)
	item, err := s.repository.getPrinterByID(s.storage, printerID, false)
	if err != nil {
		return printerEntity{}, hueat_err.ErrGeneric
	}
	if hueat_utils.IsEmpty(item) {
		return printerEntity{}, errPrinterNotFound
	}
	return item, nil
}

func (s printerService) createPrinter(ctx *gin.Context, input createPrinterInputDto) (printerEntity, error) {
	now := time.Now()
	newPrinter := printerEntity{
		ID:        uuid.New(),
		Title:     input.Title,
		Url:       input.Url,
		Active:    hueat_utils.BoolPtr(false),
		Inside:    hueat_utils.BoolPtr(true),
		Outside:   hueat_utils.BoolPtr(true),
		CreatedAt: now,
		UpdatedAt: now,
	}
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		if item, err := s.repository.getPrinterByTitle(tx, input.Title, false); err != nil {
			return hueat_err.ErrGeneric
		} else if !hueat_utils.IsEmpty(item) {
			return errPrinterSameTitleAlreadyExists
		} else if _, err = s.repository.savePrinter(tx, newPrinter, hueat_db.Create); err != nil {
			return hueat_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicPrinterV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.PrinterCreatedEvent,
				EventEntity: &hueat_pubsub.PrinterEventEntity{
					ID:        newPrinter.ID,
					Title:     newPrinter.Title,
					Url:       newPrinter.Url,
					Active:    newPrinter.Active,
					Inside:    newPrinter.Inside,
					Outside:   newPrinter.Outside,
					CreatedAt: newPrinter.CreatedAt,
					UpdatedAt: newPrinter.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(printerEntity{}, newPrinter),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return printerEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return newPrinter, nil
}

func (s printerService) updatePrinter(ctx *gin.Context, input updatePrinterInputDto) (printerEntity, error) {
	now := time.Now()
	var updatedPrinter printerEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		printerId := uuid.MustParse(input.ID)
		currentPrinter, err := s.repository.getPrinterByID(tx, printerId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentPrinter) {
			return errPrinterNotFound
		}
		updatedPrinter = currentPrinter
		updatedPrinter.UpdatedAt = now
		// If the input contains a new title, check for collision
		if input.Title != nil {
			printerSameTitle, err := s.repository.getPrinterByTitle(tx, *input.Title, false)
			if err != nil {
				return hueat_err.ErrGeneric
			}
			if !hueat_utils.IsEmpty(printerSameTitle) && printerSameTitle.ID.String() != input.ID {
				return errPrinterSameTitleAlreadyExists
			}
			updatedPrinter.Title = *input.Title
		}
		if input.Active != nil {
			updatedPrinter.Active = input.Active
		}
		if input.Inside != nil {
			updatedPrinter.Inside = input.Inside
		}
		if input.Outside != nil {
			updatedPrinter.Outside = input.Outside
		}
		if input.Url != nil {
			updatedPrinter.Url = *input.Url
		}
		if _, err = s.repository.savePrinter(tx, updatedPrinter, hueat_db.Update); err != nil {
			return hueat_err.ErrGeneric
		}

		// Send an event of printer updated
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicPrinterV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.PrinterUpdatedEvent,
				EventEntity: &hueat_pubsub.PrinterEventEntity{
					ID:        updatedPrinter.ID,
					Title:     updatedPrinter.Title,
					Url:       updatedPrinter.Url,
					Active:    updatedPrinter.Active,
					Inside:    updatedPrinter.Inside,
					Outside:   updatedPrinter.Outside,
					CreatedAt: updatedPrinter.CreatedAt,
					UpdatedAt: updatedPrinter.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentPrinter, updatedPrinter),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return printerEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return updatedPrinter, nil
}

func (s printerService) deletePrinter(ctx *gin.Context, input deletePrinterInputDto) (printerEntity, error) {
	var currentPrinter printerEntity
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		// Check if exists
		printerId := uuid.MustParse(input.ID)
		currentPrinter, err := s.repository.getPrinterByID(tx, printerId, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		if hueat_utils.IsEmpty(currentPrinter) {
			return errPrinterNotFound
		}
		s.repository.deletePrinter(tx, currentPrinter)

		// Send an event of printer deleted
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicPrinterV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.PrinterDeletedEvent,
				EventEntity: &hueat_pubsub.PrinterEventEntity{
					ID:        currentPrinter.ID,
					Title:     currentPrinter.Title,
					Url:       currentPrinter.Url,
					Active:    currentPrinter.Active,
					Inside:    currentPrinter.Inside,
					Outside:   currentPrinter.Outside,
					CreatedAt: currentPrinter.CreatedAt,
					UpdatedAt: currentPrinter.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(currentPrinter, printerEntity{}),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return printerEntity{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return currentPrinter, nil
}
