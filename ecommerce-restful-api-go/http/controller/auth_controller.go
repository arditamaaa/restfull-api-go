package controller

import (
	"simple-api-go/http/model"
	"simple-api-go/http/request"
	"simple-api-go/http/response"
	"simple-api-go/http/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService  service.AuthService
	UserService  service.UserService
	TokenService service.TokenService
}

func NewAuthController(
	authService service.AuthService, userService service.UserService,
	tokenService service.TokenService,
) *AuthController {
	return &AuthController{
		AuthService:  authService,
		UserService:  userService,
		TokenService: tokenService,
	}
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	req := new(request.Register)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Register(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.Success[model.User]{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Register successfully",
			Data:    *user,
		})
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	req := new(request.Login)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Login(c, req)
	if err != nil {
		return err
	}

	token, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			User:    *user,
			Token:   token,
		})
}

func (a *AuthController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}
	if err := a.AuthService.Logout(c, token); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Logout successfully",
		})
}
