package auth

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_timeout"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type authRouterInterface interface {
	register(engine *gin.RouterGroup)
}

type authRouter struct {
	service authServiceInterface
}

func newAuthRouter(service authServiceInterface) authRouter {
	return authRouter{
		service: service,
	}
}

// Implementation
func (r authRouter) register(router *gin.RouterGroup) {
	router.POST(
		"/auth/login",
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input
			var request loginUserInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.login(ctx, request.Username, request.Password)
			if err == errInvalidCredentials {
				hueat_router.ReturnBadRequestError(ctx, err)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "auth-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/auth/refresh",
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			var request refreshTokenInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			item, err := r.service.refreshToken(ctx, request.RefreshToken)
			if err == errExpiredRefreshToken {
				hueat_router.ReturnUnauthorizedError(ctx)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "auth-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnOk(ctx, &gin.H{"item": item})
		})

	router.POST(
		"/auth/logout",
		hueat_timeout.TimeoutMiddleware(time.Duration(1)*time.Second),
		func(ctx *gin.Context) {
			// Input
			var request refreshTokenInputDto
			if err := hueat_router.BindParameters(ctx, &request); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			if err := request.validate(); err != nil {
				hueat_router.ReturnValidationError(ctx, err)
				return
			}
			// Business Logic
			err := r.service.revokeRefreshToken(ctx, request.RefreshToken)
			if err == errExpiredRefreshToken {
				hueat_router.ReturnUnauthorizedError(ctx)
				return
			}
			// Errors and output handler
			if err != nil {
				zap.L().Error("Something went wrong", zap.String("service", "auth-router"), zap.Error(err))
				hueat_router.ReturnGenericError(ctx)
				return
			}
			hueat_router.ReturnNoContent(ctx)
		})
}
