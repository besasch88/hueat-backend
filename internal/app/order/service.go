package order

import (
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"gorm.io/gorm"
)

type orderServiceInterface interface {
	createOrderFromEvent(event hueat_pubsub.TableEventEntity) error
	getOrder(ctx *gin.Context, input getOrderInputDto) (orderEntityWithChilds, error)
	updateOrder(ctx *gin.Context, input updateOrderInputDto) (orderEntityWithChilds, error)
}

type orderService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  orderRepositoryInterface
}

func newOrderService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository orderRepositoryInterface) orderService {
	return orderService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s orderService) createOrderFromEvent(event hueat_pubsub.TableEventEntity) error {
	now := time.Now()
	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		orderID := uuid.New()
		newOrder := orderEntity{
			ID:        orderID,
			TableID:   event.ID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.repository.saveOrder(s.storage, newOrder, hueat_db.Create)
		// Send an event of order created
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicOrderV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.OrderCreatedEvent,
				EventEntity: &hueat_pubsub.OrderEventEntity{
					ID:        newOrder.ID,
					TableID:   newOrder.TableID,
					CreatedAt: newOrder.CreatedAt,
					UpdatedAt: newOrder.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(orderEntity{}, newOrder),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		newCourse := courseEntity{
			ID:        uuid.New(),
			OrderID:   orderID,
			Position:  1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.repository.saveCourse(s.storage, newCourse, hueat_db.Create)
		// Send an event of course created
		if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicCourseV1, hueat_pubsub.PubSubMessage{
			Message: hueat_pubsub.PubSubEvent{
				EventID:   uuid.New(),
				EventTime: time.Now(),
				EventType: hueat_pubsub.CourseCreatedEvent,
				EventEntity: &hueat_pubsub.CourseEventEntity{
					ID:        newCourse.ID,
					OrderID:   newCourse.OrderID,
					Position:  newCourse.Position,
					CreatedAt: newCourse.CreatedAt,
					UpdatedAt: newCourse.UpdatedAt,
				},
				EventChangedFields: hueat_utils.DiffStructs(courseEntity{}, newCourse),
			},
		}); err != nil {
			return err
		} else {
			eventsToPublish = append(eventsToPublish, event)
		}
		return nil
	})
	if errTransaction != nil {
		return errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return nil
}

func (s orderService) getOrder(ctx *gin.Context, input getOrderInputDto) (orderEntityWithChilds, error) {
	tableID := uuid.MustParse(input.TableID)
	requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
	userID := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
	// Check if the user has a permission
	if slices.Contains(requester.Permissions, hueat_auth.READ_OTHER_TABLES) {
		userID = nil
	}
	// Check if the table exists for that user
	if item, err := s.repository.checkTableExists(s.storage, userID, tableID); err != nil {
		return orderEntityWithChilds{}, hueat_err.ErrGeneric
	} else if hueat_utils.IsEmpty(item) {
		return orderEntityWithChilds{}, errTableNotFound
	}
	// Retrieve the order
	order, err := s.repository.getOrderByTableID(s.storage, tableID, false)
	if err != nil {
		return orderEntityWithChilds{}, hueat_err.ErrGeneric
	}
	if hueat_utils.IsEmpty(order) {
		return orderEntityWithChilds{}, errOrderNotFound
	}
	// Retrieve all courses
	courses, _, err := s.repository.listCoursesByOrderID(s.storage, order.ID, false)
	if err != nil {
		return orderEntityWithChilds{}, hueat_err.ErrGeneric
	}
	// For each course, retrieve its selections
	coursesWithChilds := []courseEntityWithChilds{}
	for _, course := range courses {
		selections, _, err := s.repository.listCourseSelectionsByCourseID(s.storage, course.ID, false)
		if err != nil {
			return orderEntityWithChilds{}, hueat_err.ErrGeneric
		}
		coursesWithChilds = append(coursesWithChilds, courseEntityWithChilds{
			courseEntity: course,
			Items:        selections,
		})
	}
	return orderEntityWithChilds{
		orderEntity: order,
		Courses:     coursesWithChilds,
	}, nil
}

func (s orderService) updateOrder(ctx *gin.Context, input updateOrderInputDto) (orderEntityWithChilds, error) {

	eventsToPublish := []hueat_pubsub.EventToPublish{}
	errTransaction := s.storage.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		tableID := uuid.MustParse(input.TableID)
		requester := hueat_auth.GetAuthenticatedUserFromSession(ctx)
		userID := hueat_utils.GetOptionalUUIDFromString(&requester.ID)
		// Check if the user has a permission
		if slices.Contains(requester.Permissions, hueat_auth.READ_OTHER_TABLES) {
			userID = nil
		}
		// Check if the table exists for that user
		if item, err := s.repository.checkTableExists(s.storage, userID, tableID); err != nil {
			return hueat_err.ErrGeneric
		} else if hueat_utils.IsEmpty(item) {
			return errTableNotFound
		}
		// Check if the Order exists
		order, err := s.repository.getOrderByTableID(s.storage, tableID, true)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		// If it does not exists, return an error
		if hueat_utils.IsEmpty(order) {
			return errOrderNotFound
		}

		// Get existing courses
		courses, total, err := s.repository.listCoursesByOrderID(s.storage, order.ID, false)
		if err != nil {
			return hueat_err.ErrGeneric
		}
		// Now for each course received from input, check if already exists or needed to be created a new one
		var lastPosition int64 = total
		for index, inputCourse := range input.Courses {
			var course courseEntity
			if len(courses) >= index+1 {
				course = courses[index]
				// Check if the order of items are respected
				if course.ID.String() != inputCourse.ID {
					return errCourseMismatch
				}
			} else {
				// If there are new course input, let's create them
				lastPosition++
				course = courseEntity{
					ID:        hueat_utils.GetUUIDFromString(inputCourse.ID),
					OrderID:   order.ID,
					Position:  lastPosition,
					CreatedAt: now,
					UpdatedAt: now,
				}
				if _, err := s.repository.saveCourse(s.storage, course, hueat_db.Create); err != nil {
					return hueat_err.ErrGeneric
				}
				// Send an event of course created
				if event, err := s.pubSubAgent.Persist(tx, hueat_pubsub.TopicCourseV1, hueat_pubsub.PubSubMessage{
					Message: hueat_pubsub.PubSubEvent{
						EventID:   uuid.New(),
						EventTime: time.Now(),
						EventType: hueat_pubsub.CourseCreatedEvent,
						EventEntity: &hueat_pubsub.CourseEventEntity{
							ID:        course.ID,
							OrderID:   course.OrderID,
							Position:  course.Position,
							CreatedAt: course.CreatedAt,
							UpdatedAt: course.UpdatedAt,
						},
						EventChangedFields: hueat_utils.DiffStructs(courseEntity{}, course),
					},
				}); err != nil {
					return hueat_err.ErrGeneric
				} else {
					eventsToPublish = append(eventsToPublish, event)
				}
			}
			if err := s.repository.deleteSelectionsByCourseID(s.storage, course.ID); err != nil {
				return hueat_err.ErrGeneric
			}
			for _, inputSelection := range inputCourse.Items {
				selection := courseSelectionEntity{
					ID:           uuid.New(),
					CourseID:     course.ID,
					MenuItemID:   hueat_utils.GetUUIDFromString(inputSelection.MenuItemID),
					MenuOptionID: hueat_utils.GetOptionalUUIDFromString(inputSelection.MenuOptionID),
					Quantity:     inputSelection.Quantity,
					Note:         inputSelection.Note,
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				if _, err := s.repository.saveSelection(s.storage, selection, hueat_db.Create); err != nil {
					return hueat_err.ErrGeneric
				}
			}
		}
		return nil
	})
	if errTransaction != nil {
		return orderEntityWithChilds{}, errTransaction
	} else {
		s.pubSubAgent.PublishBulk(eventsToPublish)
	}
	return s.getOrder(ctx, getOrderInputDto{TableID: input.TableID})
}
