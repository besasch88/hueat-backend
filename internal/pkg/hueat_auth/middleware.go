package hueat_auth

import (
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_utils"
	"github.com/gin-gonic/gin"
)

var authConfig AuthConfig

type AuthConfig struct {
	JwtSecret string
}

/*
Initialize the AuthMiddleware by setting the JWT secret needed to validate the received token
*/
func InitAuthMiddleware(config AuthConfig) {
	authConfig = config
}

/*
AuthMiddleware Middleware on APIs to check if the user is authenticated
and verify the permissions the user has compared to the permissions required
by the API. In case of failure, returns an error to the client.
*/
func AuthMiddleware(permissionsToCheck []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the authenticated user
		authenticatedUser, err := getAuthenticatedUserFromRequest(ctx)
		// In case of error or if the user is not found, return Unauthorized
		if err != nil || hueat_utils.IsEmpty(authenticatedUser) {
			hueat_router.ReturnUnauthorizedError(ctx)
			return
		}
		// If permissions are not defined, return Forbidden by default.
		if len(permissionsToCheck) == 0 {
			hueat_router.ReturnForbiddenError(ctx)
			return
		}
		// Check if all the required permissions are included in the authenticated User permissions, otherwise return Forbidden
		if !containsAll(authenticatedUser.Permissions, permissionsToCheck) {
			hueat_router.ReturnForbiddenError(ctx)
			return
		}
		ctx.Set(contextAuthenticatedUser, &authenticatedUser)
		ctx.Next()
	}
}

/*
GetAuthenticatedUserFromSession retrieves the authenticated user from the session.
This works in combination of the Authentication middleware that extracts all the information
provided by the JWT sent in the Authentication header of the request and store them
in the request context. This utility retrieve the authenticated user from the context session
without performing additional read operations to get all the users information.
*/
func GetAuthenticatedUserFromSession(ctx *gin.Context) *AuthenticatedUser {
	value, exists := ctx.Get(contextAuthenticatedUser)
	if exists {
		return value.(*AuthenticatedUser)
	}
	return nil
}
