package middleware

import (
	"simple-api-go/http/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Role(t service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization, Please Login")
		}

		userToken, err := t.GetTokenWithUser(c, token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}
		if userToken.User == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate, User not found")
		}
		if userToken.User.Role != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "You don't have permission to access this resource")
		}

		return c.Next()
	}
}
