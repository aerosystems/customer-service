package handlers

import (
	"encoding/json"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomerHandler struct {
	customerUsecase CustomerUsecase
}

func NewCustomerHandler(
	customerUsecase CustomerUsecase,
) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

type CreateCustomerRequest struct {
	CreateCustomerRequestBody
}

type CreateCustomerRequestBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
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
// @Success 201 {object} handlers.Customer
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/customers [post]
func (ch CustomerHandler) CreateCustomer(c echo.Context) error {
	var req CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "could not bind request")
	}

	var customerReq Customer
	if err := json.Unmarshal(req.Message.Data, &customerReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not unmarshal request")
	}

	customer, err := ch.customerUsecase.CreateCustomer(customerReq.Uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ModelToCustomerResponse(customer))
}
