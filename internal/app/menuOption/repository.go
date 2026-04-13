package menuOption

import (
	"fmt"

	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuOptionRepositoryInterface interface {
	checkMenuItemExists(tx *gorm.DB, menuItemID uuid.UUID) (bool, error)
	listMenuOptions(tx *gorm.DB, menuItemID uuid.UUID, forUpdate bool) ([]menuOptionEntity, int64, error)
	getMenuOptionByID(tx *gorm.DB, menuOptionID uuid.UUID, forUpdate bool) (menuOptionEntity, error)
	getMenuOptionByTitle(tx *gorm.DB, menuOptionTitle string, forUpdate bool) (menuOptionEntity, error)
	saveMenuOption(tx *gorm.DB, menuOption menuOptionEntity, operation hueat_db.SaveOperation) (menuOptionEntity, error)
	deleteMenuOption(tx *gorm.DB, menuOption menuOptionEntity) (menuOptionEntity, error)
	recalculateMenuOptionsPosition(tx *gorm.DB, menuItemID uuid.UUID) ([]menuOptionEntity, error)
}

type menuOptionRepository struct {
	relevanceThresholdConfig float64
}

func newMenuOptionRepository(relevanceThresholdConfig float64) menuOptionRepository {
	return menuOptionRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r menuOptionRepository) checkMenuItemExists(tx *gorm.DB, menuItemID uuid.UUID) (bool, error) {
	var model *menuItemModel
	query := tx.Where("id = ?", menuItemID)
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 || hueat_utils.IsEmpty(model) {
		return false, nil
	}
	return true, nil
}

func (r menuOptionRepository) listMenuOptions(tx *gorm.DB, menuItemID uuid.UUID, forUpdate bool) ([]menuOptionEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*menuOptionModel
	query := tx.Model(menuOptionModel{}).Where("menu_item_id = ?", menuItemID)
	queryCount := tx.Model(menuOptionModel{}).Where("menu_item_id = ?", menuItemID)

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

func (r menuOptionRepository) getMenuOptionByID(tx *gorm.DB, menuOptionID uuid.UUID, forUpdate bool) (menuOptionEntity, error) {
	var model *menuOptionModel
	query := tx.Where("id = ?", menuOptionID)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuOptionEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuOptionEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuOptionRepository) getMenuOptionByTitle(tx *gorm.DB, menuOptionTitle string, forUpdate bool) (menuOptionEntity, error) {
	var model *menuOptionModel
	query := tx.Where("title = ?", menuOptionTitle)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return menuOptionEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return menuOptionEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r menuOptionRepository) saveMenuOption(tx *gorm.DB, menuOption menuOptionEntity, operation hueat_db.SaveOperation) (menuOptionEntity, error) {
	var model = menuOptionModel(menuOption)
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
		return menuOptionEntity{}, err
	}
	return menuOption, nil
}

func (r menuOptionRepository) deleteMenuOption(tx *gorm.DB, menuOption menuOptionEntity) (menuOptionEntity, error) {
	var model = menuOptionModel(menuOption)
	err := tx.Delete(model).Error
	if err != nil {
		return menuOptionEntity{}, err
	}
	return menuOption, nil
}

func (r menuOptionRepository) recalculateMenuOptionsPosition(tx *gorm.DB, menuItemID uuid.UUID) ([]menuOptionEntity, error) {
	var models []*menuOptionModel
	err := tx.Raw(`
		WITH ordered AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY position ASC, updated_at DESC) AS new_position
			FROM hueat_menu_option
			WHERE menu_item_id = ?
		)
		UPDATE hueat_menu_option s
		SET position = o.new_position, updated_at = NOW()
		FROM ordered o
		WHERE s.id = o.id
		AND s.position IS DISTINCT FROM o.new_position
		RETURNING s.*;
	`, menuItemID).Scan(&models).Error
	if err != nil {
		return []menuOptionEntity{}, err
	}
	// Return only updated entities
	var entities []menuOptionEntity = []menuOptionEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}
