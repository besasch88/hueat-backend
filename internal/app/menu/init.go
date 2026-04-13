package menu

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
	zap.L().Info("Initialize Menu package...")
	var repository menuRepositoryInterface
	var service menuServiceInterface
	var router menuRouterInterface

	repository = newMenuRepository(envs.SearchRelevanceThreshold)
	service = newMenuService(dbStorage, pubSubAgent, repository)
	router = newMenuRouter(service)
	router.register(routerGroup)
	zap.L().Info("Menu package initialized")
}
