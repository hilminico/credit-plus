package main

import (
	"creditPlus/helper/localization"
	"creditPlus/helper/validation"
	"creditPlus/internal/config"
	"creditPlus/internal/interface/handler"
	"creditPlus/internal/repository"
	"creditPlus/internal/usecase"
	"creditPlus/migration"
	"creditPlus/seeder"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// multi language
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
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
	})

	migration.RunMigration()

	seeder.RunSeeder()

	// Initialize database
	db := config.GetDB()

	// Initialize repository, service, and controller user
	customerRepo := repository.NewCustomerRepository(db)
	customerService := usecase.NewCustomerService(customerRepo)
	customerController := handler.NewCustomerController(customerService)

	// Validation
	validation.InitValidator()

	// localization
	if err := localization.InitLocalization(); err != nil {
		e.Logger.Fatal("Failed to initialize localization: ", err)
	}

	// Routes
	api := e.Group("/api/v1")
	{
		api.POST("/customer/login", customerController.Login)

		//// user route
		//api.POST("/users", userController.CreateUser)
		//api.GET("/users", userController.GetAllUsers)
		//api.GET("/users/:id", userController.GetUser)
		//api.PUT("/users/:id", userController.UpdateUser)
		//
		//// user route
		//api.POST("/login", authController.Login)
		//api.POST("/register", authController.Register)
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
