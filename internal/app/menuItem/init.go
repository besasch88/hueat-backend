package menuItem

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
	zap.L().Info("Initialize MenuItem package...")
	var repository menuItemRepositoryInterface
	var service menuItemServiceInterface
	var consumer menuItemConsumerInterface
	var router menuItemRouterInterface

	repository = newMenuItemRepository(envs.SearchRelevanceThreshold)
	service = newMenuItemService(dbStorage, pubSubAgent, repository)
	consumer = newMenuItemConsumer(pubSubAgent, service)
	consumer.subscribe()
	router = newMenuItemRouter(service)
	router.register(routerGroup)
	zap.L().Info("MenuItem package initialized")
}
