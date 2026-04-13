package statistics

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type statisticsRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type statisticsRouter struct {
	service statisticsServiceInterface
}

func newStatisticsRouter(service statisticsServiceInterface) statisticsRouter {
	return statisticsRouter{
		service: service,
	}
}

// Implementation
func (r statisticsRouter) register(router *gin.RouterGroup) {

	router.GET(
		"/statistics",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_STATISTICS}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Business Logic
			item, err := r.service.getStatistics(ctx)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "statistics-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

}
