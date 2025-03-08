package config

import (
	"encoding/json"
	"simple-api-go/util"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	return fiber.Config{
		ServerHeader:  "Fiber",
		AppName:       "Fiber API (Ecommers-App)",
		ErrorHandler:  util.ErrorHandler,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		CaseSensitive: true,
	}
}
