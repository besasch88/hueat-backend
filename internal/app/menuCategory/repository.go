package menuCategory

import (
	"fmt"

	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuCategoryRepositoryInterface interface {
	checkPrinterExists(tx *gorm.DB, printerID uuid.UUID) (bool, error)
	listMenuCategories(tx *gorm.DB, forUpdate bool) ([]menuCategoryEntity, int64, error)
	getMenuCategoryByID(tx *gorm.DB, menuCategoryID uuid.UUID, forUpdate bool) (menuCategoryEntity, error)
	getMenuCategoryByTitle(tx *gorm.DB, menuCategoryTitle string, forUpdate bool) (menuCategoryEntity, error)
	saveMenuCategory(tx *gorm.DB, menuCategory menuCategoryEntity, operation hueat_db.SaveOperation) (menuCategoryEntity, error)
	deleteMenuCategory(tx *gorm.DB, menuCategory menuCategoryEntity) (menuCategoryEntity, error)
	recalculateMenuCategorysPosition(tx *gorm.DB) ([]menuCategoryEntity, error)
}

type menuCategoryRepository struct {
	relevanceThresholdConfig float64
}

func newMenuCategoryRepository(relevanceThresholdConfig float64) menuCategoryRepository {
	return menuCategoryRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r menuCategoryRepository) checkPrinterExists(tx *gorm.DB, printerID uuid.UUID) (bool, error) {
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

func (r menuCategoryRepository) listMenuCategories(tx *gorm.DB, forUpdate bool) ([]menuCategoryEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuCategoryModel
	query := tx.Model(menuCategoryModel{})
	queryCount := tx.Model(menuCategoryModel{})
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

func (r menuCategoryRepository) getMenuCategoryByID(tx *gorm.DB, menuCategoryID uuid.UUID, forUpdate bool) (menuCategoryEntity, error) {
	var model *menuCategoryModel
	query := tx.Where("id = ?", menuCategoryID)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuCategoryEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuCategoryEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuCategoryRepository) getMenuCategoryByTitle(tx *gorm.DB, menuCategoryTitle string, forUpdate bool) (menuCategoryEntity, error) {
	var model *menuCategoryModel
	query := tx.Where("title = ?", menuCategoryTitle)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuCategoryEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuCategoryEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuCategoryRepository) saveMenuCategory(tx *gorm.DB, menuCategory menuCategoryEntity, operation hueat_db.SaveOperation) (menuCategoryEntity, error) {
	var model = menuCategoryModel(menuCategory)
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
		return menuCategoryEntity{}, err
	}
	return menuCategory, nil
}

func (r menuCategoryRepository) deleteMenuCategory(tx *gorm.DB, menuCategory menuCategoryEntity) (menuCategoryEntity, error) {
	var model = menuCategoryModel(menuCategory)
	err := tx.Delete(model).Error
	if err != nil {
		return menuCategoryEntity{}, err
	}
	return menuCategory, nil
}

func (r menuCategoryRepository) recalculateMenuCategorysPosition(tx *gorm.DB) ([]menuCategoryEntity, error) {
	var models []*menuCategoryModel
	err := tx.Raw(`
		WITH ordered AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY position ASC, updated_at DESC) AS new_position
			FROM hueat_menu_category
		)
		UPDATE hueat_menu_category t
		SET position = o.new_position, updated_at = NOW()
		FROM ordered o
		WHERE t.id = o.id
		AND t.position IS DISTINCT FROM o.new_position
		RETURNING t.*;
	`).Scan(&models).Error
	if err != nil {
		return []menuCategoryEntity{}, err
	}
	// Return only updated entities
	var entities []menuCategoryEntity = []menuCategoryEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}
