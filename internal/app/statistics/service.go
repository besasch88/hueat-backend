package statistics

import (
	"github.com/hueat/backend/internal/pkg/hueat_err"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type statisticsServiceInterface interface {
	getStatistics(ctx *gin.Context) (statisticsEntity, error)
}

type statisticsService struct {
	storage     *gorm.DB
	pubSubAgent *hueat_pubsub.PubSubAgent
	repository  statisticsRepositoryInterface
}

func newStatisticsService(storage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, repository statisticsRepositoryInterface) statisticsService {
	return statisticsService{
		storage:     storage,
		pubSubAgent: pubSubAgent,
		repository:  repository,
	}
}

func (s statisticsService) getStatistics(ctx *gin.Context) (statisticsEntity, error) {
	avgTableDuration, err := s.repository.getAverageTableDuration(s.storage)
	if err != nil {
		return statisticsEntity{}, hueat_err.ErrGeneric
	}
	paymentMethodsTakins, err := s.repository.getPaymentMethodsTakins(s.storage)
	if err != nil {
		return statisticsEntity{}, hueat_err.ErrGeneric
	}
	menuItemStats, err := s.repository.getMenuItemStats(s.storage)
	if err != nil {
		return statisticsEntity{}, hueat_err.ErrGeneric
	}
	return statisticsEntity{
		AvgTableDuration: avgTableDuration,
		PaymentsTakins:   paymentMethodsTakins,
		MenuItemStats:    menuItemStats,
	}, nil
}
