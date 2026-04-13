package table

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type tableRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type tableRouter struct {
	service tableServiceInterface
}

func newTableRouter(service tableServiceInterface) tableRouter {
	return tableRouter{
		service: service,
	}
}

// Implementation
func (r tableRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/tables",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request listTablesInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			items, totalCount, err := r.service.listTables(ctx, request)

			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/tables",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES, hueat_auth.WRITE_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createTableInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createTable(ctx, request)
			if err == errTableSameNameAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/tables/:tableId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getTableInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getTableByID(ctx, request)
			if err == errTableNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/tables/:tableId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES, hueat_auth.WRITE_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updateTableInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updateTable(ctx, request)
			if err == errTableNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errTableSameNameAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/tables/:tableId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_MY_TABLES, hueat_auth.WRITE_MY_TABLES}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deleteTableInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deleteTable(ctx, request)
			if err == errTableNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "table-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnNoContent(ctx)
		})

}
