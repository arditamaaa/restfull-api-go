package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerConfig() fiber.Handler {
	logConfig := logger.Config{
		Format:     "${time} ${method} ${status} ${path} | ${latency}\n",
		TimeFormat: "15:04:05.00",
	}
	logger := logger.New(logConfig)

	// if _, err := os.Stat("logs"); os.IsNotExist(err) {
	// 	err := os.Mkdir("logs", 0770)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// timeStamp := time.Now().Format("2006-01-02_03:04:05")
	// file, err := os.OpenFile("logs/log_"+timeStamp+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer file.Close()
	// multiWriter := io.MultiWriter(os.Stdout, file)
	// logConfig.Output = multiWriter
	return logger
}
