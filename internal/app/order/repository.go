package order

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepositoryInterface interface {
	checkTableExists(tx *gorm.DB, userID *uuid.UUID, tableID uuid.UUID) (bool, error)
	getOrderByTableID(tx *gorm.DB, tableID uuid.UUID, forUpdate bool) (orderEntity, error)
	listCoursesByOrderID(tx *gorm.DB, orderID uuid.UUID, forUpdate bool) ([]courseEntity, int64, error)
	listCourseSelectionsByCourseID(tx *gorm.DB, courseID uuid.UUID, forUpdate bool) ([]courseSelectionEntity, int64, error)
	saveOrder(tx *gorm.DB, order orderEntity, operation hueat_db.SaveOperation) (orderEntity, error)
	saveCourse(tx *gorm.DB, course courseEntity, operation hueat_db.SaveOperation) (courseEntity, error)
	saveSelection(tx *gorm.DB, selection courseSelectionEntity, operation hueat_db.SaveOperation) (courseSelectionEntity, error)
	deleteSelectionsByCourseID(tx *gorm.DB, courseID uuid.UUID) error
	getOrderDetailByTableID(tx *gorm.DB, tableID uuid.UUID) ([]orderDetailEntity, error)
	getOrderDetailByTableIDAndCourseID(tx *gorm.DB, tableID uuid.UUID, courseID uuid.UUID) ([]orderDetailEntity, error)
	getPricedOrderByTableID(tx *gorm.DB, tableID uuid.UUID) ([]orderDetailEntity, error)
	getTotalPriceAndPaymentByTableID(tx *gorm.DB, tableID uuid.UUID) (paymentDetailEntity, error)
}

type orderRepository struct {
	relevanceThresholdConfig float64
}

func newOrderRepository(relevanceThresholdConfig float64) orderRepository {
	return orderRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r orderRepository) checkTableExists(tx *gorm.DB, userID *uuid.UUID, tableID uuid.UUID) (bool, error) {
	var model *tableModel
	query := tx.Where("id = ?", tableID)
	if userID != nil {
		query.Where("user_id = ?", userID)
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 || hueat_utils.IsEmpty(model) {
		return false, nil
	}
	return true, nil
}

func (r orderRepository) getOrderByTableID(tx *gorm.DB, tableID uuid.UUID, forUpdate bool) (orderEntity, error) {
	var model *orderModel
	query := tx.Where("table_id = ?", tableID)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return orderEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return orderEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r orderRepository) listCoursesByOrderID(tx *gorm.DB, orderID uuid.UUID, forUpdate bool) ([]courseEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*courseModel
	query := tx.Model(courseModel{}).Where("order_id = ?", orderID)
	queryCount := tx.Model(courseModel{}).Where("order_id = ?", orderID)

	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", hueat_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []courseEntity{}, 0, result.Error
	}
	var entities []courseEntity = []courseEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r orderRepository) listCourseSelectionsByCourseID(tx *gorm.DB, courseID uuid.UUID, forUpdate bool) ([]courseSelectionEntity, int64, error) {
	var totalCount int64

	var models []*courseSelectionModel
	query := tx.Model(courseSelectionModel{}).Where("course_id = ?", courseID)
	queryCount := tx.Model(courseSelectionModel{}).Where("course_id = ?", courseID)

	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []courseSelectionEntity{}, 0, result.Error
	}
	var entities []courseSelectionEntity = []courseSelectionEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r orderRepository) saveOrder(tx *gorm.DB, order orderEntity, operation hueat_db.SaveOperation) (orderEntity, error) {
	var model = orderModel(order)
	var err error
	switch operation {
	case hueat_db.Create:
		err = tx.Create(model).Error
	case hueat_db.Update:
		err = tx.Updates(model).Error
	case hueat_db.Upsert:
		err = tx.Save(model).Error
	}
	if err != nil {
		return orderEntity{}, err
	}
	return order, nil
}

func (r orderRepository) saveCourse(tx *gorm.DB, course courseEntity, operation hueat_db.SaveOperation) (courseEntity, error) {
	var model = courseModel(course)
	var err error
	switch operation {
	case hueat_db.Create:
		err = tx.Create(model).Error
	case hueat_db.Update:
		err = tx.Updates(model).Error
	case hueat_db.Upsert:
		err = tx.Save(model).Error
	}
	if err != nil {
		return courseEntity{}, err
	}
	return course, nil
}

func (r orderRepository) saveSelection(tx *gorm.DB, selection courseSelectionEntity, operation hueat_db.SaveOperation) (courseSelectionEntity, error) {
	var model = courseSelectionModel(selection)
	var err error
	switch operation {
	case hueat_db.Create:
		err = tx.Create(model).Error
	case hueat_db.Update:
		err = tx.Updates(model).Error
	case hueat_db.Upsert:
		err = tx.Save(model).Error
	}
	if err != nil {
		return courseSelectionEntity{}, err
	}
	return selection, nil
}

func (r orderRepository) deleteSelectionsByCourseID(tx *gorm.DB, courseID uuid.UUID) error {
	err := tx.Where("course_id = ?", courseID).Delete(&courseSelectionModel{}).Error
	return err
}

func (r orderRepository) getOrderDetailByTableID(tx *gorm.DB, tableID uuid.UUID) ([]orderDetailEntity, error) {
	var models []orderDetailModel
	err := tx.Raw(getOrderByTableQuery, tableID).Scan(&models).Error
	if err != nil {
		return []orderDetailEntity{}, err
	}
	var entities []orderDetailEntity = []orderDetailEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}
func (r orderRepository) getOrderDetailByTableIDAndCourseID(tx *gorm.DB, tableID uuid.UUID, courseID uuid.UUID) ([]orderDetailEntity, error) {
	var models []orderDetailModel
	err := tx.Raw(getCourseByTableAndCourseQuery, tableID, courseID).Scan(&models).Error
	if err != nil {
		return []orderDetailEntity{}, err
	}
	var entities []orderDetailEntity = []orderDetailEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}

func (r orderRepository) getPricedOrderByTableID(tx *gorm.DB, tableID uuid.UUID) ([]orderDetailEntity, error) {
	var models []orderDetailModel
	err := tx.Raw(getPricedOrderByTableQuery, tableID).Scan(&models).Error
	if err != nil {
		return []orderDetailEntity{}, err
	}
	var entities []orderDetailEntity = []orderDetailEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}

func (r orderRepository) getTotalPriceAndPaymentByTableID(tx *gorm.DB, tableID uuid.UUID) (paymentDetailEntity, error) {
	var model paymentDetailModel
	err := tx.Raw(getTotalPriceAndPaymentByTableQuery, tableID).First(&model).Error
	if err == gorm.ErrRecordNotFound {
		return paymentDetailEntity{}, nil
	} else if err != nil {
		return paymentDetailEntity{}, err
	}
	return model.toEntity(), nil
}
