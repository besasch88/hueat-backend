package menuOption

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuOptionRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuOptionRouter struct {
	service menuOptionServiceInterface
}

func newMenuOptionRouter(service menuOptionServiceInterface) menuOptionRouter {
	return menuOptionRouter{
		service: service,
	}
}

// Implementation
func (r menuOptionRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/menu/items/:menuItemId/options",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request listMenuOptionsInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			items, totalCount, err := r.service.listMenuOptions(ctx, request)
			if err == errMenuItemNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/menu/items/:menuItemId/options",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createMenuOptionInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createMenuOption(ctx, request)
			if err == errMenuItemNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuOptionSameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/menu/options/:menuOptionId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuOptionInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenuOptionByID(ctx, request)
			if err == errMenuOptionNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/menu/options/:menuOptionId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateMenuOptionInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateMenuOption(ctx, request)
			if err == errMenuOptionNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuOptionSameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/menu/options/:menuOptionId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteMenuOptionInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteMenuOption(ctx, request)
			if err == errMenuOptionNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-option-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnNoContent(ctx)
		})
}
