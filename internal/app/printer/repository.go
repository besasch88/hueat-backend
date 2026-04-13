package printer

import (
	"fmt"

	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type printerRepositoryInterface interface {
	listPrinters(tx *gorm.DB, forUpdate bool) ([]printerEntity, int64, error)
	getPrinterByID(tx *gorm.DB, printerID uuid.UUID, forUpdate bool) (printerEntity, error)
	getPrinterByTitle(tx *gorm.DB, printerTitle string, forUpdate bool) (printerEntity, error)
	savePrinter(tx *gorm.DB, printer printerEntity, operation hueat_db.SaveOperation) (printerEntity, error)
	deletePrinter(tx *gorm.DB, printer printerEntity) (printerEntity, error)
}

type printerRepository struct {
	relevanceThresholdConfig float64
}

func newPrinterRepository(relevanceThresholdConfig float64) printerRepository {
	return printerRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r printerRepository) listPrinters(tx *gorm.DB, forUpdate bool) ([]printerEntity, int64, error) {
	var totalCount int64
	var order string

	var models []*printerModel
	query := tx.Model(printerModel{})
	queryCount := tx.Model(printerModel{})
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	order = fmt.Sprintf("%s %s", "created_at", hueat_db.Asc)
	result := query.Order(order).Find(&models)
	queryCount.Count(&totalCount)

	if result.Error != nil {
		return []printerEntity{}, 0, result.Error
	}
	var entities []printerEntity = []printerEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, totalCount, nil
}

func (r printerRepository) getPrinterByID(tx *gorm.DB, printerID uuid.UUID, forUpdate bool) (printerEntity, error) {
	var model *printerModel
	query := tx.Where("id = ?", printerID)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return printerEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return printerEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r printerRepository) getPrinterByTitle(tx *gorm.DB, printerTitle string, forUpdate bool) (printerEntity, error) {
	var model *printerModel
	query := tx.Where("title = ?", printerTitle)
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return printerEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return printerEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r printerRepository) savePrinter(tx *gorm.DB, printer printerEntity, operation hueat_db.SaveOperation) (printerEntity, error) {
	var model = printerModel(printer)
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
		return printerEntity{}, err
	}
	return printer, nil
}

func (r printerRepository) deletePrinter(tx *gorm.DB, printer printerEntity) (printerEntity, error) {
	var model = printerModel(printer)
	err := tx.Delete(model).Error
	if err != nil {
		return printerEntity{}, err
	}
	return printer, nil
}
