package util

import (
	"simple-api-go/http/model"

	"github.com/gofiber/fiber/v2"
)

func AuthUser(c *fiber.Ctx) model.User {
	user := model.User{}
	authUser := c.Locals("auth_user")
	if authUser != nil {
		user = *authUser.(*model.User)
	}
	return user
}

func AuthUserId(c *fiber.Ctx) uint64 {
	user := AuthUser(c)
	return user.ID
}
