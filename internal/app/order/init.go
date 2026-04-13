package order

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
	zap.L().Info("Initialize Order package...")
	var repository orderRepositoryInterface
	var printRepository printRepositoryInterface
	var service orderServiceInterface
	var printService printServiceInterface
	var consumer orderConsumerInterface
	var router orderRouterInterface

	repository = newOrderRepository(envs.SearchRelevanceThreshold)
	printRepository = newPrintRepository()
	service = newOrderService(dbStorage, pubSubAgent, repository)
	printService = newPrintService(envs.PrinterEnabled, dbStorage, pubSubAgent, repository, printRepository)
	consumer = newOrderConsumer(pubSubAgent, service)
	consumer.subscribe()
	router = newOrderRouter(service, printService)
	router.register(routerGroup)
	zap.L().Info("Order package initialized")
}
