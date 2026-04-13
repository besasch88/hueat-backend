package statistics

import (
	"github.com/hueat/backend/internal/pkg/hueat_env"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
Init the module by registering new APIs and PubSub consumers.
*/
func Init(envs *hueat_env.Envs, dbStorage *gorm.DB, pubSubAgent *hueat_pubsub.PubSubAgent, routerGroup *gin.RouterGroup) {
	zap.L().Info("Initialize Statistics package...")
	var repository statisticsRepositoryInterface
	var service statisticsServiceInterface
	var router statisticsRouterInterface

	repository = newStatisticsRepository(envs.SearchRelevanceThreshold)
	service = newStatisticsService(dbStorage, pubSubAgent, repository)
	router = newStatisticsRouter(service)
	router.register(routerGroup)
	zap.L().Info("Statistics package initialized")
}
