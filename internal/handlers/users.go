package handlers

import (
	"github.com/aerosystems/user-service/internal/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetUser godoc
// @Summary Get user
// @Description Get user
// @Tags users
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} Response{data=models.User}
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /v1/users [get]
func (h *BaseHandler) GetUser(c echo.Context) error {
	accessTokenClaims := c.Get("accessTokenClaims").(services.AccessTokenClaims)
	user, err := h.userService.GetUserById(uint(accessTokenClaims.UserId))
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not get user", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "user was successfully found", user)
}
