package middlewares

import (
	"creditPlus/helper/localization"
	"github.com/labstack/echo/v4"
)

func WithLocalization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// set default id lang
			lang := "en"

			if c.Request().Header.Get("Accept-Language") != "" {
				lang = c.Request().Header.Get("Accept-Language")
			}

			ctx := localization.WithLanguage(c.Request().Context(), lang)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
