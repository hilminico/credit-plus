package middlewares

import (
	"creditPlus/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

type AuthConfig struct {
	DB          *gorm.DB
	SigningKey  []byte
	ContextKey  string
	TokenLookup string
}

func DefaultAuthConfig(db *gorm.DB) AuthConfig {
	return AuthConfig{
		DB:          db,
		SigningKey:  []byte(os.Getenv("SECRET_KEY")),
		ContextKey:  "customer",
		TokenLookup: "header:Authorization",
	}
}

func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token
			tokenString, err := extractToken(c, config.TokenLookup)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			// Parse token
			token, err := jwt.ParseWithClaims(tokenString, &domain.CustomerClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
				}
				return config.SigningKey, nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Validate claims
			if claims, ok := token.Claims.(*domain.CustomerClaims); ok && token.Valid {
				// Check if user exists in database
				var customer domain.Customer
				if err := config.DB.Where("unique_identifier = ?", claims.UniqueIdentifier).First(&customer).Error; err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "customer not found")
				}

				// Store user in context
				c.Set(config.ContextKey, &customer)
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}
	}
}

func extractToken(c echo.Context, lookup string) (string, error) {
	parts := strings.Split(lookup, ":")
	if len(parts) != 2 {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "invalid token lookup format")
	}

	switch parts[0] {
	case "header":
		auth := c.Request().Header.Get(parts[1])
		if auth == "" {
			return "", echo.NewHTTPError(http.StatusBadRequest, "missing auth header")
		}
		return strings.TrimPrefix(auth, "Bearer "), nil
	default:
		return "", echo.NewHTTPError(http.StatusInternalServerError, "unsupported token lookup")
	}
}
