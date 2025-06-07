package HTTPServer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/aerosystems/customer-service/internal/entities"
)

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
// @Param raw body CreateFirebaseCustomerRequestBody true "Create user"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/firebase/create-customer [post]
func (h Handler) CreateCustomer(c echo.Context) error {
	var req CreateFirebaseCustomerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not bind request")
	}

	var customerReq CreateCustomerFirebaseEvent
	if err := json.Unmarshal(req.Message.Data, &customerReq); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "could not unmarshal request")
	}

	err := h.customerUsecase.CreateCustomer(c.Request().Context(), customerReq.Email, customerReq.UID)
	if errors.Is(err, entities.ErrCustomerAlreadyExists) {
		c.NoContent(http.StatusCreated)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}
