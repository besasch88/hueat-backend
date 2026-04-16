package menu

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuRepositoryInterface interface {
	getTableByID(tx *gorm.DB, tableID uuid.UUID) (tableEntity, error)
	listMenuCategories(tx *gorm.DB, inside bool, forUpdate bool) ([]menuCategoryEntity, int64, error)
	listMenuItems(tx *gorm.DB, inside bool, targetTableID uuid.UUID, forUpdate bool) ([]menuItemEntity, int64, error)
	listMenuOptions(tx *gorm.DB, inside bool, forUpdate bool) ([]menuOptionEntity, int64, error)
}

type menuRepository struct {
	relevanceThresholdConfig float64
}

func newMenuRepository(relevanceThresholdConfig float64) menuRepository {
	return menuRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r menuRepository) getTableByID(tx *gorm.DB, tableID uuid.UUID) (tableEntity, error) {
	var model *tableModel
	query := tx.Where("id = ?", tableID)
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return tableEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return tableEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuRepository) listMenuCategories(tx *gorm.DB, inside bool, forUpdate bool) ([]menuCategoryEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuCategoryModel
	query := tx.Model(menuCategoryModel{})
	queryCount := tx.Model(menuCategoryModel{})
	if inside {
		query.Where("inside = ?", true)
		queryCount.Where("inside = ?", true)
	} else {
		query.Where("outside = ?", true)
		queryCount.Where("outside = ?", true)
	}
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", hueat_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []menuCategoryEntity{}, 0, result.Error
	}
	var entities []menuCategoryEntity = []menuCategoryEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r menuRepository) listMenuItems(tx *gorm.DB, inside bool, targetTableID uuid.UUID, forUpdate bool) ([]menuItemEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuItemModel
	query := tx.Model(menuItemModel{})
	queryCount := tx.Model(menuItemModel{})
	if inside {
		query.Where("inside = ?", true)
		queryCount.Where("inside = ?", true)
	} else {
		query.Where("outside = ?", true)
		queryCount.Where("outside = ?", true)
	}
	query.Where(tx.Where("table_id IS NULL").Or("table_id = ?", targetTableID))
	queryCount.Where(tx.Where("table_id IS NULL").Or("table_id = ?", targetTableID))
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", hueat_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []menuItemEntity{}, 0, result.Error
	}
	var entities []menuItemEntity = []menuItemEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r menuRepository) listMenuOptions(tx *gorm.DB, inside bool, forUpdate bool) ([]menuOptionEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuOptionModel
	query := tx.Model(menuOptionModel{})
	queryCount := tx.Model(menuOptionModel{})
	if inside {
		query.Where("inside = ?", true)
		queryCount.Where("inside = ?", true)
	} else {
		query.Where("outside = ?", true)
		queryCount.Where("outside = ?", true)
	}
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "position", hueat_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []menuOptionEntity{}, 0, result.Error
	}
	var entities []menuOptionEntity = []menuOptionEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}
