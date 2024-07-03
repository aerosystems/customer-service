package CustomErrors

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EchoError struct {
	mode   EchoHandlerMode
	errors []ApiError
}

func NewEchoErrorHandler(mode string) echo.HTTPErrorHandler {
	e := EchoError{
		mode:   NewEchoHandlerMode(mode),
		errors: list,
	}
	return e.Handler
}

func (h *EchoError) Handler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message interface{}

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		if httpError.Internal != nil {
			var herr *echo.HTTPError
			if errors.As(httpError.Internal, &herr) {
				httpError = herr
			}
		}
	}

	var apiError *ApiError
	if errors.As(err, &apiError) {
		httpError = &echo.HTTPError{
			Code:    apiError.HttpCode,
			Message: apiError.Message,
		}
	}

	code = httpError.Code
	message = httpError.Message

	if _, ok := httpError.Message.(string); ok {
		message = map[string]interface{}{"message": err.Error()}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(httpError.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}

func unwrapError(err error) error {
	originalError := err

	for originalError != nil {
		internalError := errors.Unwrap(originalError)
		if internalError == nil {
			break
		}
		originalError = internalError
	}
	return originalError
}
