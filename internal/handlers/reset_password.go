package handlers

import (
	"errors"
	"github.com/aerosystems/auth-service/pkg/normalizers"
	"github.com/aerosystems/auth-service/pkg/validators"
	"net/http"
)

type ResetPasswordRequestBody struct {
	Email    string `json:"email" example:"example@gmail.com"`
	Password string `json:"password" example:"P@ssw0rd"`
}

// ResetPassword godoc
// @Summary resetting password
// @Description Password should contain:
// @Description - minimum of one small case letter
// @Description - minimum of one upper case letter
// @Description - minimum of one digit
// @Description - minimum of one special character
// @Description - minimum 8 characters length
// @Tags auth
// @Accept  json
// @Produce application/json
// @Param registration body handlers.ResetPasswordRequestBody true "raw request body"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/user/reset-password [post]
func (h *BaseHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var requestPayload ResetPasswordRequestBody
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "could not read request body", err))
		return
	}
	addr, err := validators.ValidateEmail(requestPayload.Email)
	if err != nil {
		err = errors.New("email is not valid")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(422005, "Email does not valid", err))
		return
	}
	email := normalizers.NormalizeEmail(addr)
	err = validators.ValidatePassword(requestPayload.Password)
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(422006, "Password does not valid", err))
		return
	}
	if err := h.userService.ResetPassword(email, requestPayload.Password); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500001, "could not reset password", err))
		return
	}
	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("password was successfully reset", nil))
	return
}
