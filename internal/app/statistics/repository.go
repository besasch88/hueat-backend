package statistics

import (
	"time"

	"gorm.io/gorm"
)

type statisticsRepositoryInterface interface {
	getAverageTableDuration(tx *gorm.DB) (time.Duration, error)
	getPaymentMethodsTakins(tx *gorm.DB) ([]paymentTakingsEntity, error)
	getMenuItemStats(tx *gorm.DB) ([]menuItemStatEntity, error)
	deleteStatistics(tx *gorm.DB) error
}

type statisticsRepository struct {
	relevanceThresholdConfig float64
}

func newStatisticsRepository(relevanceThresholdConfig float64) statisticsRepository {
	return statisticsRepository{
		relevanceThresholdConfig: relevanceThresholdConfig,
	}
}

func (r statisticsRepository) getAverageTableDuration(tx *gorm.DB) (time.Duration, error) {
	var model averageTableDurationModel
	err := tx.Raw(getAverageTableDurationQuery).First(&model).Error
	if err == gorm.ErrRecordNotFound {
		return time.Duration(0), nil
	} else if err != nil {
		return time.Duration(0), err
	}
	return model.Duration, nil
}

func (r statisticsRepository) getPaymentMethodsTakins(tx *gorm.DB) ([]paymentTakingsEntity, error) {
	var models []paymentTakingsModel
	err := tx.Raw(getPaymentMethodsTakinsQuery).Scan(&models).Error
	if err != nil {
		return []paymentTakingsEntity{}, err
	}
	var entities []paymentTakingsEntity = []paymentTakingsEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}

func (r statisticsRepository) getMenuItemStats(tx *gorm.DB) ([]menuItemStatEntity, error) {
	var models []menuItemStatModel
	err := tx.Raw(getMenuItemStatsQuery).Scan(&models).Error
	if err != nil {
		return []menuItemStatEntity{}, err
	}
	var entities []menuItemStatEntity = []menuItemStatEntity{}
	for _, model := range models {
		entity := model.toEntity()
		entities = append(entities, entity)
	}
	return entities, nil
}

func (r statisticsRepository) deleteStatistics(tx *gorm.DB) error {
	return tx.Model(&tableModel{}).Where("close IS TRUE").Delete(&tableModel{}).Error
}
