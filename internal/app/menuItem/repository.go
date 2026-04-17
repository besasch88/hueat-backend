package menuItem

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuItemRepositoryInterface interface {
	getTableByID(tx *gorm.DB, userID *uuid.UUID, tableID uuid.UUID) (tableEntity, error)
	checkPrinterExists(tx *gorm.DB, printerID uuid.UUID) (bool, error)
	checkMenuCategoryExists(tx *gorm.DB, menuCategoryID uuid.UUID) (bool, error)
	listMenuItems(tx *gorm.DB, menuCategoryID uuid.UUID, forUpdate bool) ([]menuItemEntity, int64, error)
	getMenuItemByID(tx *gorm.DB, menuItemID uuid.UUID, forUpdate bool) (menuItemEntity, error)
	getCustomMenuItemByID(tx *gorm.DB, menuItemID uuid.UUID, forUpdate bool) (menuItemEntity, error)
	getMenuItemByTitle(tx *gorm.DB, menuItemTitle string, forUpdate bool) (menuItemEntity, error)
	saveMenuItem(tx *gorm.DB, menuItem menuItemEntity, operation hueat_db.SaveOperation) (menuItemEntity, error)
	deleteMenuItem(tx *gorm.DB, menuItem menuItemEntity) (menuItemEntity, error)
	recalculateMenuItemsPosition(tx *gorm.DB, menuCategoryID uuid.UUID) ([]menuItemEntity, error)
}

type menuItemRepository struct {
	relevanceThresholdConfig float64
}

func newMenuItemRepository(relevanceThresholdConfig float64) menuItemRepository {
	return menuItemRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r menuItemRepository) getTableByID(tx *gorm.DB, userID *uuid.UUID, tableID uuid.UUID) (tableEntity, error) {
	var model *tableModel
	query := tx.Where("id = ?", tableID).Where("close IS FALSE")
	if userID != nil {
		query.Where("user_id = ?", userID)
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

func (r menuItemRepository) checkPrinterExists(tx *gorm.DB, printerID uuid.UUID) (bool, error) {
	var model *printerModel
	query := tx.Where("id = ?", printerID)
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 || hueat_utils.IsEmpty(model) {
		return false, nil
	}
	return true, nil
}

func (r menuItemRepository) checkMenuCategoryExists(tx *gorm.DB, menuCategoryID uuid.UUID) (bool, error) {
	var model *menuCategoryModel
	query := tx.Where("id = ?", menuCategoryID)
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 || hueat_utils.IsEmpty(model) {
		return false, nil
	}
	return true, nil
}

func (r menuItemRepository) listMenuItems(tx *gorm.DB, menuCategoryID uuid.UUID, forUpdate bool) ([]menuItemEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuItemModel
	query := tx.Model(menuItemModel{}).Where("menu_category_id = ?", menuCategoryID).Where("table_id IS NULL")
	queryCount := tx.Model(menuItemModel{}).Where("menu_category_id = ?", menuCategoryID).Where("table_id IS NULL")

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

func (r menuItemRepository) getMenuItemByID(tx *gorm.DB, menuItemID uuid.UUID, forUpdate bool) (menuItemEntity, error) {
	var model *menuItemModel
	query := tx.Where("id = ?", menuItemID).Where("table_id IS NULL")
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuItemEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuItemEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuItemRepository) getCustomMenuItemByID(tx *gorm.DB, menuItemID uuid.UUID, forUpdate bool) (menuItemEntity, error) {
	var model *menuItemModel
	query := tx.Where("id = ?", menuItemID)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuItemEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuItemEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuItemRepository) getMenuItemByTitle(tx *gorm.DB, menuItemTitle string, forUpdate bool) (menuItemEntity, error) {
	var model *menuItemModel
	query := tx.Where("title = ?", menuItemTitle).Where("table_id IS NULL")
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuItemEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuItemEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuItemRepository) saveMenuItem(tx *gorm.DB, menuItem menuItemEntity, operation hueat_db.SaveOperation) (menuItemEntity, error) {
	var model = menuItemModel(menuItem)
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
		return menuItemEntity{}, err
	}
	return menuItem, nil
}

func (r menuItemRepository) deleteMenuItem(tx *gorm.DB, menuItem menuItemEntity) (menuItemEntity, error) {
	var model = menuItemModel(menuItem)
	err := tx.Delete(model).Error
	if err != nil {
		return menuItemEntity{}, err
	}
	return menuItem, nil
}

func (r menuItemRepository) recalculateMenuItemsPosition(tx *gorm.DB, menuCategoryID uuid.UUID) ([]menuItemEntity, error) {
	var models []*menuItemModel
	err := tx.Raw(`
		WITH ordered AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY position ASC, updated_at DESC) AS new_position
			FROM hueat_menu_item
			WHERE menu_category_id = ?
		)
		UPDATE hueat_menu_item s
		SET position = o.new_position, updated_at = NOW()
		FROM ordered o
		WHERE s.id = o.id
		AND s.position IS DISTINCT FROM o.new_position
		RETURNING s.*;
	`, menuCategoryID).Scan(&models).Error
	if err != nil {
		return []menuItemEntity{}, err
	}
	// Return only updated entities
	var entities []menuItemEntity = []menuItemEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}
