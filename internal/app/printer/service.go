package printer

import (
	"time"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_db"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_err"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_utils"
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
	pubSubAgent *ceng_pubsub.PubSubAgent
	repository  printerRepositoryInterface
}

func newPrinterService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository printerRepositoryInterface) printerService {
	return printerService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s printerService) listPrinters(ctx *gin.Context) ([]printerEntity, int64, error) {
	items, totalCount, err := s.repository.listPrinters(s.storage, false)
	if err != nil || items == nil {
		return []printerEntity{}, 0, ceng_err.ErrGeneric
	}
	return items, totalCount, nil
}

func (s printerService) getPrinterByID(ctx *gin.Context, input getPrinterInputDto) (printerEntity, error) {
	printerID := uuid.MustParse(input.ID)
	item, err := s.repository.getPrinterByID(s.storage, printerID, false)
	if err != nil {
		return printerEntity{}, ceng_err.ErrGeneric
	}
	if ceng_utils.IsEmpty(item) {
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
		Active:    ceng_utils.BoolPtr(false),
		Inside:    ceng_utils.BoolPtr(true),
		Outside:   ceng_utils.BoolPtr(true),
		CreatedAt: now,
		UpdatedAt: now,
	}
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		if item, err := s.repository.getPrinterByTitle(tx, input.Title, false); err != nil {
			return ceng_err.ErrGeneric
		} else if !ceng_utils.IsEmpty(item) {
			return errPrinterSameTitleAlreadyExists
		} else if _, err = s.repository.savePrinter(tx, newPrinter, ceng_db.Create); err != nil {
			return ceng_err.ErrGeneric
		}

		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicPrinterV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.PrinterCreatedEvent,
				EventEntity: &ceng_pubsub.PrinterEventEntity{
					ID:        newPrinter.ID,
					Title:     newPrinter.Title,
					Url:       newPrinter.Url,
					Active:    newPrinter.Active,
					Inside:    newPrinter.Inside,
					Outside:   newPrinter.Outside,
					CreatedAt: newPrinter.CreatedAt,
					UpdatedAt: newPrinter.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(printerEntity{}, newPrinter),
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
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		printerId := uuid.MustParse(input.ID)
		currentPrinter, err := s.repository.getPrinterByID(tx, printerId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentPrinter) {
			return errPrinterNotFound
		}
		updatedPrinter = currentPrinter
		updatedPrinter.UpdatedAt = now
		// If the input contains a new title, check for collision
		if input.Title != nil {
			printerSameTitle, err := s.repository.getPrinterByTitle(tx, *input.Title, false)
			if err != nil {
				return ceng_err.ErrGeneric
			}
			if !ceng_utils.IsEmpty(printerSameTitle) && printerSameTitle.ID.String() != input.ID {
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
		if _, err = s.repository.savePrinter(tx, updatedPrinter, ceng_db.Update); err != nil {
			return ceng_err.ErrGeneric
		}

		// Send an event of printer updated
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicPrinterV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.PrinterUpdatedEvent,
				EventEntity: &ceng_pubsub.PrinterEventEntity{
					ID:        updatedPrinter.ID,
					Title:     updatedPrinter.Title,
					Url:       updatedPrinter.Url,
					Active:    updatedPrinter.Active,
					Inside:    updatedPrinter.Inside,
					Outside:   updatedPrinter.Outside,
					CreatedAt: updatedPrinter.CreatedAt,
					UpdatedAt: updatedPrinter.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentPrinter, updatedPrinter),
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
	eventsToPublish := []ceng_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		// Check if exists
		printerId := uuid.MustParse(input.ID)
		currentPrinter, err := s.repository.getPrinterByID(tx, printerId, true)
		if err != nil {
			return ceng_err.ErrGeneric
		}
		if ceng_utils.IsEmpty(currentPrinter) {
			return errPrinterNotFound
		}
		s.repository.deletePrinter(tx, currentPrinter)

		// Send an event of printer deleted
		if event, err := s.pubSubAgent.Persist(tx, ceng_pubsub.TopicPrinterV1, ceng_pubsub.PubSubMessage{
			Message: ceng_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: ceng_pubsub.PrinterDeletedEvent,
				EventEntity: &ceng_pubsub.PrinterEventEntity{
					ID:        currentPrinter.ID,
					Title:     currentPrinter.Title,
					Url:       currentPrinter.Url,
					Active:    currentPrinter.Active,
					Inside:    currentPrinter.Inside,
					Outside:   currentPrinter.Outside,
					CreatedAt: currentPrinter.CreatedAt,
					UpdatedAt: currentPrinter.UpdatedAt,
				},
				EventChangedFields: ceng_utils.DiffStructs(currentPrinter, printerEntity{}),
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
