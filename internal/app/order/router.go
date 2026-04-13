package order

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type orderRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type orderRouter struct {
	service      orderServiceInterface
	printService printServiceInterface
}

func newOrderRouter(service orderServiceInterface, printService printServiceInterface) orderRouter {
	return orderRouter{
		service:      service,
		printService: printService,
	}
}

// Implementation
func (r orderRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/tables/:tableId/order",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getOrderInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getOrder(ctx, request)
			if err == errTableNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errOrderNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "order-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/tables/:tableId/order",
		hueat_auth.AuthMiddleware([]string{hueat_auth.WRITE_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateOrderInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateOrder(ctx, request)
			if err == errTableNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errOrderNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errCourseMismatch {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}

			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "order-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/tables/:tableId/order/print",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES, hueat_auth.WRITE_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(5)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request printOrderInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			err := r.printService.print(ctx, request)
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "order-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"success": true})
		})

}
