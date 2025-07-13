package response

import (
	"creditPlus/helper/localization"
	"creditPlus/helper/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(c echo.Context, status int, messageID string, templateData interface{}) error {
	message := localization.Localize(c.Request().Context(), messageID, map[string]interface{}{})
	return c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    templateData,
	})
}

func ErrorResponse(c echo.Context, status int, messageID string, templateData map[string]interface{}) error {
	message := localization.Localize(c.Request().Context(), messageID, templateData)
	return c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

func ErrorResponseValidation(c echo.Context, errors []validation.ValidationError) error {
	return c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Bad Request",
		Errors:  errors,
	})
}
