package middleware

import (
	"simple-api-go/http/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(t service.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		method := c.Method()
		path := c.Route().Path
		url := strings.TrimPrefix(path, "/api")

		if token == "" && url == "/carts/add" && method == "POST" {
			return fiber.NewError(fiber.StatusUnauthorized, "Please login first to add items to your cart")
		}

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

		c.Locals("auth_user", userToken.User)
		return c.Next()
	}
}
