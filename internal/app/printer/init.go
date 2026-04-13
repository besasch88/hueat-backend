package printer

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
	zap.L().Info("Initialize Printer package...")
	var repository printerRepositoryInterface
	var service printerServiceInterface
	var router printerRouterInterface

	repository = newPrinterRepository(envs.SearchRelevanceThreshold)
	service = newPrinterService(dbStorage, pubSubAgent, repository)
	router = newPrinterRouter(service)
	router.register(routerGroup)
	zap.L().Info("Printer package initialized")
}
