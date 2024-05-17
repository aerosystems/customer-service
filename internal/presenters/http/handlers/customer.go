package handlers

import (
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomerHandler struct {
	*BaseHandler
	customerUsecase CustomerUsecase
}

func NewCustomerHandler(
	baseHandler *BaseHandler,
	customerUsecase CustomerUsecase,
) *CustomerHandler {
	return &CustomerHandler{
		BaseHandler:     baseHandler,
		customerUsecase: customerUsecase,
	}
}

type CreateCustomerRequest struct {
	CreateCustomerRequestBody
}

type CreateCustomerRequestBody struct {
	Customer
}

type Customer struct {
	Uuid string `json:"uuid"`
}

func ModelToCustomerResponse(user *models.Customer) *Customer {
	return &Customer{
		Uuid: user.Uuid.String(),
	}
}

// CreateCustomer godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param raw body handlers.CreateCustomerRequestBody true "Create user"
// @Success 201 {object} handlers.Response{data=handlers.Customer}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 401 {object} handlers.ErrorResponse
// @Failure 403 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /v1/customers [post]
func (ch CustomerHandler) CreateCustomer(c echo.Context) error {
	var req CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return ch.ErrorResponse(c, http.StatusBadRequest, "could not bind request", err)
	}
	user, err := ch.customerUsecase.CreateCustomer(req.Uuid)
	if err != nil {
		if customErr, ok := err.(*CustomErrors.ConflictError); ok {
			return ch.ErrorResponse(c, http.StatusConflict, err.Error(), customErr)
		}
		return ch.ErrorResponse(c, http.StatusInternalServerError, "could not create user", err)
	}
	return ch.SuccessResponse(c, http.StatusCreated, "customer was successfully created", ModelToCustomerResponse(user))
}
