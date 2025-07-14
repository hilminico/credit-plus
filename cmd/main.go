package main

import (
	"creditPlus/helper/localization"
	"creditPlus/helper/validation"
	"creditPlus/internal/config"
	"creditPlus/internal/interface/handler"
	"creditPlus/internal/repository"
	"creditPlus/internal/usecase"
	"creditPlus/middlewares"
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

	// multi language fetch from header
	e.Use(middlewares.WithLocalization())

	migration.RunMigration()

	seeder.RunSeeder()

	// Initialize database
	db := config.GetDB()

	// Initialize repository, service, and controller user
	customerRepo := repository.NewCustomerRepository(db)
	customerDetailRepo := repository.NewCustomerDetailRepository(db)
	customerService := usecase.NewCustomerService(customerRepo, customerDetailRepo)
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
		api.POST("/login", customerController.Login)

		// Restricted group
		restricted := api.Group("/customer")

		// Configure middleware with the custom claims type
		restricted.Use(middlewares.AuthWithConfig(middlewares.DefaultAuthConfig(db)))

		restricted.GET("/profile", customerController.Show)
		restricted.PATCH("/", customerController.Update)

	}

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
