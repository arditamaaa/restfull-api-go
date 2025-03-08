package util

import (
	"errors"
	res "simple-api-go/http/response"
	"simple-api-go/http/validation"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if errorsMap := validation.CustomErrorMessages(err); len(errorsMap) > 0 {
		return res.Error(c, fiber.StatusBadRequest, "Bad Request", errorsMap)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return res.Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	return res.Error(c, fiber.StatusInternalServerError, "Internal Server Error", nil)
}

func NotFoundHandler(c *fiber.Ctx) error {
	return res.Error(c, fiber.StatusNotFound, "Endpoint Not Found", nil)
}
