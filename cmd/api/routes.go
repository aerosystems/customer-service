package main

import (
	_ "github.com/aerosystems/customer-service/docs" // docs are generated by Swag CLI, you have to import it.
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (app *Config) NewRouter() *echo.Echo {
	e := echo.New()

	docsGroup := e.Group("/docs")
	docsGroup.Use(app.basicAuthMiddleware.BasicAuthMiddleware)
	docsGroup.GET("/*", echoSwagger.WrapHandler)

	e.GET("/v1/customers", app.baseHandler.GetCustomer, app.oauthMiddleware.AuthTokenMiddleware(models.CustomerRole))

	return e
}
