package menu

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuRouter struct {
	service menuServiceInterface
}

func newMenuRouter(service menuServiceInterface) menuRouter {
	return menuRouter{
		service: service,
	}
}

// Implementation
func (r menuRouter) register(router *gin.RouterGroup) {
	router.GET(
		"tables/:tableId/menu",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenu(ctx, request)
			// Errors and output handler
			if err == errTableNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

}
