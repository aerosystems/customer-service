package HttpServer

import (
	OAuthService "github.com/aerosystems/customer-service/pkg/oauth"
)

type TokenService interface {
	GetAccessSecret() string
	DecodeAccessToken(tokenString string) (*OAuthService.AccessTokenClaims, error)
}
