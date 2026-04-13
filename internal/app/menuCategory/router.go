package menuCategory

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type menuCategoryRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type menuCategoryRouter struct {
	service menuCategoryServiceInterface
}

func newMenuCategoryRouter(service menuCategoryServiceInterface) menuCategoryRouter {
	return menuCategoryRouter{
		service: service,
	}
}

// Implementation
func (r menuCategoryRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/menu/categories",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Business Logic
			items, totalCount, err := r.service.listMenuCategories(ctx)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/menu/categories",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createMenuCategoryInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createMenuCategory(ctx, request)
			if err == errMenuCategorySameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			if err == errPrinterNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/menu/categories/:menuCategoryId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getMenuCategoryInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getMenuCategoryByID(ctx, request)
			if err == errMenuCategoryNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/menu/categories/:menuCategoryId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateMenuCategoryInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateMenuCategory(ctx, request)
			if err == errMenuCategoryNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errMenuCategorySameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/menu/categories/:menuCategoryId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MENU, hueat_auth.WRITE_MENU}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteMenuCategoryInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteMenuCategory(ctx, request)
			if err == errMenuCategoryNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "menu-category-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnNoContent(ctx)
		})
}
