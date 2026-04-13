package table

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
	zap.L().Info("Initialize Table package...")
	var repository tableRepositoryInterface
	var service tableServiceInterface
	var router tableRouterInterface

	repository = newTableRepository(envs.SearchRelevanceThreshold)
	service = newTableService(dbStorage, pubSubAgent, repository)
	router = newTableRouter(service)
	router.register(routerGroup)
	zap.L().Info("Table package initialized")
}
