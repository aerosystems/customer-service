package handlers

import (
	"encoding/json"
	"errors"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type FirebaseHandler struct {
	customerUsecase CustomerUsecase
}

func NewFirebaseHandler(
	customerUsecase CustomerUsecase,
) *FirebaseHandler {
	return &FirebaseHandler{
		customerUsecase: customerUsecase,
	}
}

type CreateFirebaseCustomerRequest struct {
	CreateFirebaseCustomerRequestBody
}

type CreateFirebaseCustomerRequestBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type CreateCustomerFirebaseEvent struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
}

// CreateCustomer godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept  json
// @Produce application/json
// @Param raw body handlers.CreateFirebaseCustomerRequestBody true "Create user"
// @Success 201 null
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/firebase/create-customer [post]
func (ch FirebaseHandler) CreateCustomer(c echo.Context) error {
	var req CreateFirebaseCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not bind request")
	}

	var customerReq CreateCustomerFirebaseEvent
	if err := json.Unmarshal(req.Message.Data, &customerReq); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "could not unmarshal request")
	}

	err := ch.customerUsecase.CreateCustomer(c.Request().Context(), customerReq.Email, customerReq.UID)
	if errors.Is(err, CustomErrors.ErrCustomerAlreadyExists) {
		c.NoContent(http.StatusCreated)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
