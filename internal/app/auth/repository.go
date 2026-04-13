package auth

import (
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type authRepositoryInterface interface {
	getAuthSessionEntityByRefreshToken(tx *gorm.DB, refreshToken string, forUpdate bool) (authSessionEntity, error)
	saveAuthSessionEntity(tx *gorm.DB, entity authSessionEntity, operation hueat_db.SaveOperation) (authSessionEntity, error)
	deleteAuthSessionEntity(tx *gorm.DB, entity authSessionEntity) error
	cleanUpExpiredRefreshToken(tx *gorm.DB) error
}

type authRepository struct {
}

func newAuthRepository() authRepository {
	return authRepository{}
}

func (r authRepository) getAuthSessionEntityByRefreshToken(tx *gorm.DB, refreshToken string, forUpdate bool) (authSessionEntity, error) {
	var model *authSessionModel
	query := tx.Where("refresh_token = ?", refreshToken)
	query.Where("expires_at > NOW()")
	if forUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := query.Limit(1).Find(&model)
	if result.Error != nil {
		return authSessionEntity{}, result.Error
	}
	if result.RowsAffected == 0 {
		return authSessionEntity{}, nil
	}
	return model.toEntity(), nil
}

func (r authRepository) saveAuthSessionEntity(tx *gorm.DB, entity authSessionEntity, operation hueat_db.SaveOperation) (authSessionEntity, error) {
	var model = authSessionModel(entity)
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
		return authSessionEntity{}, err
	}
	return entity, nil
}

func (r authRepository) deleteAuthSessionEntity(tx *gorm.DB, entity authSessionEntity) error {
	var model = authSessionModel(entity)
	if err := tx.Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func (r authRepository) cleanUpExpiredRefreshToken(tx *gorm.DB) error {
	return tx.Where("expires_at < NOW()").Delete(&authSessionModel{}).Error
}
