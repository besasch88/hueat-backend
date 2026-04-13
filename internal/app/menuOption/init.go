package menuOption

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
	zap.L().Info("Initialize MenuOption package...")
	var repository menuOptionRepositoryInterface
	var service menuOptionServiceInterface
	var router menuOptionRouterInterface

	repository = newMenuOptionRepository(envs.SearchRelevanceThreshold)
	service = newMenuOptionService(dbStorage, pubSubAgent, repository)
	router = newMenuOptionRouter(service)
	router.register(routerGroup)
	zap.L().Info("MenuOption package initialized")
}
