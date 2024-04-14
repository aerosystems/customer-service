package handlers

import (
	"errors"
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

type CustomerResponse struct {
	Uuid string `json:"uuid"`
}

func ModelToCustomerResponse(user *models.Customer) *CustomerResponse {
	return &CustomerResponse{
		Uuid: user.Uuid.String(),
	}
}

// GetCustomer godoc
// @Summary Get user
// @Description Get user
// @Tags users
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} rest.Response{data=models.Customer}
// @Failure 401 {object} rest.ErrorResponse
// @Failure 403 {object} rest.ErrorResponse
// @Failure 404 {object} rest.ErrorResponse
// @Failure 500 {object} rest.ErrorResponse
// @Router /v1/customers [get]
func (ch CustomerHandler) GetCustomer(c echo.Context) error {
	accessTokenClaims := c.Get("accessTokenClaims").(models.AccessTokenClaims)
	user, err := ch.customerUsecase.GetUserByUuid(accessTokenClaims.UserUuid)
	if err != nil {
		return ch.ErrorResponse(c, http.StatusInternalServerError, "could not get user", err)
	}
	if user == nil {
		err := errors.New("user not found")
		return ch.ErrorResponse(c, http.StatusNotFound, err.Error(), err)
	}
	return ch.SuccessResponse(c, http.StatusOK, "customer was successfully found", ModelToCustomerResponse(user))
}
