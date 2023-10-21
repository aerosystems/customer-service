package middleware

import (
	"errors"
	"github.com/aerosystems/user-service/internal/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"net/http"
	"os"
	"strings"
	"time"
)

func AuthTokenMiddleware(roles []string) echo.MiddlewareFunc {
	AuthorizationConfig := echojwt.Config{
		SigningKey:     []byte(os.Getenv("ACCESS_SECRET")),
		ParseTokenFunc: parseToken,
		ErrorHandler: func(c echo.Context, err error) error {
			return handlers.ErrorResponse(c, http.StatusUnauthorized, "invalid token", err)
		},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return AuthorizationConfig.ErrorHandler(c, errors.New("missing Authorization header"))
			}

			// Token should be in the form "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return AuthorizationConfig.ErrorHandler(c, errors.New("invalid token format"))
			}

			token := tokenParts[1]

			accessTokenClaims, err := TokenService.DecodeAccessToken(token)
			if err != nil {
				return AuthorizationConfig.ErrorHandler(c, err)
			}

			if int64(accessTokenClaims.Exp) < time.Now().Unix() {
				return AuthorizationConfig.ErrorHandler(c, errors.New("token expired"))
			}

			// Перевірка, чи userRole знаходиться у переданому масиві roles
			roleFound := false
			for _, role := range roles {
				if accessTokenClaims.UserRole == role {
					roleFound = true
					break
				}
			}

			if !roleFound {
				// Якщо userRole не входить у roles, повертаємо помилку 403
				return handlers.ErrorResponse(c, http.StatusForbidden, "forbidden", errors.New("userRole not allowed"))
			}

			return next(c)
		}
	}
}

func parseToken(c echo.Context, auth string) (interface{}, error) {
	_ = c
	accessTokenClaims, err := TokenService.DecodeAccessToken(auth)
	if err != nil {
		return nil, err
	}

	return accessTokenClaims, nil
}
