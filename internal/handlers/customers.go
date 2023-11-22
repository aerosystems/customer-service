package handlers

import (
	"github.com/aerosystems/customer-service/internal/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetCustomer godoc
// @Summary Get user
// @Description Get user
// @Tags users
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} handlers.Response{data=models.Customer}
// @Failure 401 {object} handlers.Response
// @Failure 403 {object} handlers.Response
// @Failure 500 {object} handlers.Response
// @Router /v1/customers [get]
func (h *BaseHandler) GetCustomer(c echo.Context) error {
	accessTokenClaims := c.Get("accessTokenClaims").(services.AccessTokenClaims)
	user, err := h.customerService.GetUserByUuid(accessTokenClaims.UserUuid)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not get user", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "customer was successfully found", user)
}
