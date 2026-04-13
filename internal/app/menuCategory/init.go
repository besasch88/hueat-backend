package menuCategory

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
	zap.L().Info("Initialize MenuCategory package...")
	var repository menuCategoryRepositoryInterface
	var service menuCategoryServiceInterface
	var router menuCategoryRouterInterface

	repository = newMenuCategoryRepository(envs.SearchRelevanceThreshold)
	service = newMenuCategoryService(dbStorage, pubSubAgent, repository)
	router = newMenuCategoryRouter(service)
	router.register(routerGroup)
	zap.L().Info("MenuCategory package initialized")
}
