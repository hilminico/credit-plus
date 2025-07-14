package handler

import (
	"creditPlus/helper/response"
	"creditPlus/helper/validation"
	"creditPlus/internal/domain"
	"creditPlus/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomerController struct {
	customerService *usecase.CustomerService
}

func NewCustomerController(authService *usecase.CustomerService) *CustomerController {
	return &CustomerController{customerService: authService}
}

func (c *CustomerController) Login(ctx echo.Context) error {
	var customerRequest domain.LoginRequest

	if validationErrors := validation.ValidateRequest(ctx, &customerRequest); len(validationErrors) > 0 {
		return response.ErrorResponseValidation(ctx, validationErrors)
	}

	loginResponse, err := c.customerService.Login(
		ctx.Request().Context(),
		&domain.LoginRequest{
			Email:    customerRequest.Email,
			Password: customerRequest.Password,
		})

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(ctx, http.StatusCreated, "customer.login", loginResponse)
}

func (c *CustomerController) Show(ctx echo.Context) error {
	customer := ctx.Get("customer").(*domain.Customer)

	user, err := c.customerService.GetCustomer(ctx.Request().Context(), customer.UniqueIdentifier)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusNotFound, "customer.not_found", nil)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "customer.retrieved", user)
}

func (c *CustomerController) Update(ctx echo.Context) error {
	customer := ctx.Get("customer").(*domain.Customer)

	var updateData domain.CustomerDetailRequest
	if validationErrors := validation.ValidateRequest(ctx, &updateData); len(validationErrors) > 0 {
		return response.ErrorResponseValidation(ctx, validationErrors)
	}

	customerDetail, err := c.customerService.UpdateCustomerDetail(ctx.Request().Context(), customer.UniqueIdentifier, &updateData)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return response.SuccessResponse(ctx, http.StatusOK, "customer.updated", customerDetail)
}
