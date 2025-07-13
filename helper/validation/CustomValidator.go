package validation

import (
	"creditPlus/helper/localization"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"reflect"
	"strings"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	// Register custom validations here
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateRequest(c echo.Context, request interface{}) []ValidationError {
	if err := c.Bind(request); err != nil {
		return []ValidationError{{Field: "request", Message: "Invalid request body"}}
	}

	var errors []ValidationError
	if err := Validate.Struct(request); err != nil {
		localizer := localization.GetLocalizer(c.Request().Context())

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			message := getValidationMessage(localizer, err)
			errors = append(errors, ValidationError{
				Field:   field,
				Message: message,
			})
		}
	}
	return errors
}

func getValidationMessage(localizer *i18n.Localizer, err validator.FieldError) string {
	fieldName := err.Field()

	switch err.Tag() {
	case "required":
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "validation.required",
			TemplateData: map[string]string{
				"field": fieldName,
			},
		})

		fmt.Println(msg)
		return msg
	case "email":
		return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "validation.email"})
	case "min":
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "validation.min",
			TemplateData: map[string]string{
				"field": fieldName,
				"param": err.Param(),
			},
		})
		return msg
	case "max":
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "validation.max",
			TemplateData: map[string]string{
				"field": fieldName,
				"param": err.Param(),
			},
		})
		return msg
	default:
		return "Invalid value"
	}
}
