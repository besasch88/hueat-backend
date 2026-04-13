package auth

import (
	"time"

	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUtilInterface interface {
	generateToken(user authUserEntity) (authTokenEntity, error)
	generateHashedPassword(plainPassword string) (string, error)
	checkPassword(plainPassword string, hashedPassword string) bool
}

type authUtil struct {
	authJwtSecret               string
	authJwtAccessTokenDuration  int
	authJwtRefreshTokenDuration int
}

func newAuthUtil(authJwtSecret string, authJwtAccessTokenDuration int, authJwtRefreshTokenDuration int) authUtil {
	return authUtil{
		authJwtSecret:               authJwtSecret,
		authJwtAccessTokenDuration:  authJwtAccessTokenDuration,
		authJwtRefreshTokenDuration: authJwtRefreshTokenDuration,
	}
}

func (u authUtil) generateToken(user authUserEntity) (authTokenEntity, error) {
	// Define JWT Claims including permissions
	type CustomClaims struct {
		jwt.RegisteredClaims
		Username    string   `json:"username"`
		Permissions []string `json:"permissions"`
	}

	// Define Access and Refresh Token ID and their duration
	now := time.Now()
	accessTokenExpiresAt := now.Add(time.Duration(u.authJwtAccessTokenDuration) * time.Second)
	refreshTokenExpiresAt := now.Add(time.Duration(u.authJwtRefreshTokenDuration) * time.Second)
	accessTokenID := uuid.New()
	refreshTokenID := uuid.New()

	// Create Access Token with claims
	accessTokenClaims := CustomClaims{
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessTokenID.String(),
			Issuer:    "hueat",
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Username: user.Username,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(u.authJwtSecret))
	if err != nil {
		return authTokenEntity{}, err
	}

	// Create Refresh Token with claims
	refreshTokenClaims := CustomClaims{
		Permissions: []string{hueat_auth.UPDATE_REFRESH_TOKEN},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshTokenID.String(),
			Issuer:    "hueat",
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Username: user.Username,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(u.authJwtSecret))
	if err != nil {
		return authTokenEntity{}, err
	}

	// Return generated tokens
	return authTokenEntity{
		UserId:                user.ID.String(),
		AccessToken:           accessTokenStr,
		RefreshToken:          refreshTokenStr,
		AccessTokenID:         accessTokenID,
		RefreshTokenID:        refreshTokenID,
		AccessTokenCreatedAt:  now,
		RefreshTokenCreatedAt: now,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func (u authUtil) generateHashedPassword(plainPassword string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost, // currently 10
	)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (u authUtil) checkPassword(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
	return err == nil
}
