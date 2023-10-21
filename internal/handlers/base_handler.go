package handlers

import (
	"github.com/aerosystems/user-service/internal/services"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
)

type BaseHandler struct {
	UserService services.UserService
}

func NewBaseHandler(userService services.UserService) *BaseHandler {
	return &BaseHandler{
		UserService: userService,
	}
}

// Response is the type used for sending JSON around
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ErrResponse is the type used for sending JSON around
type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SuccessResponse takes a response status code and arbitrary data and writes a json response to the client in depends on Header Accept
func SuccessResponse(c echo.Context, statusCode int, message string, data any) error {
	payload := Response{
		Message: message,
		Data:    data,
	}
	return c.JSON(statusCode, payload)
}

// ErrorResponse takes a response status code and arbitrary data and writes a json response to the client in depends on Header Accept and APP_ENV environment variable(has two possible values: dev and prod)
// - APP_ENV=dev responds debug info level of error
// - APP_ENV=prod responds just message about error [DEFAULT]
func ErrorResponse(c echo.Context, statusCode int, message string, err error) error {
	// TODO: add custom codes for errors
	payload := ErrResponse{
		Code:    statusCode,
		Message: message,
	}

	if strings.ToLower(os.Getenv("APP_ENV")) == "dev" {
		payload.Data = err.Error()
	}

	return c.JSON(statusCode, payload)
}
