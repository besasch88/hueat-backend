package menuItem

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuItemRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuItemRouter struct {
	service menuItemServiceInterface
}

func newMenuItemRouter(service menuItemServiceInterface) menuItemRouter {
	return menuItemRouter{
		service: service,
	}
}

// Implementation
func (r menuItemRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/menu/categories/:menuCategoryId/items",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request listMenuItemsInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			items, totalCount, err := r.service.listMenuItems(ctx, request)
			if err == errMenuCategoryNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-item-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/menu/categories/:menuCategoryId/items",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createMenuItemInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createMenuItem(ctx, request)
			if err == errMenuCategoryNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuItemSameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-item-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/menu/items/:menuItemId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuItemInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenuItemByID(ctx, request)
			if err == errMenuItemNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-item-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/menu/items/:menuItemId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateMenuItemInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateMenuItem(ctx, request)
			if err == errMenuItemNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuItemSameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-item-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/menu/items/:menuItemId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteMenuItemInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteMenuItem(ctx, request)
			if err == errMenuItemNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-item-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnNoContent(ctx)
		})
}
