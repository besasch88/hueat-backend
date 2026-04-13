package printer

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type printerRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type printerRouter struct {
	service printerServiceInterface
}

func newPrinterRouter(service printerServiceInterface) printerRouter {
	return printerRouter{
		service: service,
	}
}

// Implementation
func (r printerRouter) register(router *gin.RouterGroup) {
	router.GET(
		"/printers",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_PRINTER}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Business Logic
			items, totalCount, err := r.service.listPrinters(ctx)
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"items": items, "totalCount": totalCount})
		})

	router.POST(
		"/printers",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_PRINTER, hueat_auth.WRITE_PRINTER}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request createPrinterInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.createPrinter(ctx, request)
			if err == errPrinterSameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.GET(
		"/printers/:printerId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_PRINTER}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request getPrinterInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.getPrinterByID(ctx, request)
			if err == errPrinterNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.PUT(
		"/printers/:printerId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_PRINTER, hueat_auth.WRITE_PRINTER}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request updatePrinterInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.updatePrinter(ctx, request)
			if err == errPrinterNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			if err == errPrinterSameTitleAlreadyExists {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.DELETE(
		"/printers/:printerId",
		hueat_auth.AuthMiddleware([]string{hueat_auth.READ_PRINTER, hueat_auth.WRITE_PRINTER}),
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input validation
			var request deletePrinterInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			_, err := r.service.deletePrinter(ctx, request)
			if err == errPrinterNotFound {
				hueat_router.ReturnNotFoundError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "printer-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnNoContent(ctx)
		})
}
