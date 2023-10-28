package handlers

import (
	"github.com/aerosystems/auth-service/pkg/normalizers"
	"github.com/aerosystems/auth-service/pkg/validators"
	"net/http"
)

type RegistrationRequestBody struct {
	Email    string `json:"email" example:"example@gmail.com"`
	Password string `json:"password" example:"P@ssw0rd"`
}

// Register godoc
// @Summary registration user by credentials
// @Description Password should contain:
// @Description - minimum of one small case letter
// @Description - minimum of one upper case letter
// @Description - minimum of one digit
// @Description - minimum of one special character
// @Description - minimum 8 characters length
// @Tags auth
// @Accept  json
// @Produce application/json
// @Param registration body handlers.RegistrationRequestBody true "raw request body"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/user/register [post]
func (h *BaseHandler) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload RegistrationRequestBody
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "could not read request body", err))
		return
	}
	addr, err := validators.ValidateEmail(requestPayload.Email)
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400005, "email does not valid", err))
		return
	}
	email := normalizers.NormalizeEmail(addr)
	if err := validators.ValidatePassword(requestPayload.Password); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400006, "password does not valid", err))
		return
	}
	if err := h.userService.Register(email, requestPayload.Password, r.RemoteAddr); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500005, "could not register user", err))
		return
	}
	_ = WriteResponse(w, http.StatusCreated, NewResponsePayload("user was registered successfully", nil))
	return
}
