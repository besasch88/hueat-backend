package table

import (
	"fmt"

	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tableRepositoryInterface interface {
	listTables(tx *gorm.DB, userId *uuid.UUID, inside bool, includeClosed *bool, forUpdate bool) ([]tableEntity, int64, error)
	getTableByID(tx *gorm.DB, tableID uuid.UUID, userId *uuid.UUID, forUpdate bool) (tableEntity, error)
	getOpenTableByName(tx *gorm.DB, tableName string, forUpdate bool) (tableEntity, error)
	saveTable(tx *gorm.DB, table tableEntity, operation hueat_db.SaveOperation) (tableEntity, error)
	deleteTable(tx *gorm.DB, table tableEntity) (tableEntity, error)
}

type tableRepository struct {
	relevanceThresholdConfig float64
}

func newTableRepository(relevanceThresholdConfig float64) tableRepository {
	return tableRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r tableRepository) listTables(tx *gorm.DB, userId *uuid.UUID, inside bool, includeClosed *bool, forUpdate bool) ([]tableEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*tableModel
	query := tx.Model(tableModel{}).Where("inside = ?", inside)
	queryCount := tx.Model(tableModel{}).Where("inside = ?", inside)

	if userId != nil {
		query = query.Where("user_id = ?", userId)
		queryCount = queryCount.Where("user_id = ?", userId)
	}
	if includeClosed == nil || !*includeClosed {
		query = query.Where("close = ?", false)
		queryCount = queryCount.Where("close = ?", false)
	}
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s, %s %s", "close", hueat_db.Asc, "created_at", hueat_db.Desc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []tableEntity{}, 0, result.Error
	}
	var entities []tableEntity = []tableEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r tableRepository) getTableByID(tx *gorm.DB, tableID uuid.UUID, userId *uuid.UUID, forUpdate bool) (tableEntity, error) {
	var model *tableModel
	query := tx.Where("id = ?", tableID)
	if userId != nil {
		query = query.Where("user_id = ?", userId)
	}
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return tableEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return tableEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r tableRepository) getOpenTableByName(tx *gorm.DB, tableName string, forUpdate bool) (tableEntity, error) {
	var model *tableModel
	query := tx.Where("name = ?", tableName).Where("close = ?", false)

	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return tableEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return tableEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r tableRepository) saveTable(tx *gorm.DB, table tableEntity, operation hueat_db.SaveOperation) (tableEntity, error) {
	var model = tableModel(table)
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
		return tableEntity{}, err
	}
	return table, nil
}

func (r tableRepository) deleteTable(tx *gorm.DB, table tableEntity) (tableEntity, error) {
	var model = tableModel(table)
	err := tx.Delete(model).Error
	if err != nil {
		return tableEntity{}, err
	}
	return table, nil
}
