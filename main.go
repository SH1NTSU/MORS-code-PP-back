package main

import (
	"MORS-code/convert"
	"MORS-code/message"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/message", message.Messanger)
	app.Post("/convert", convert.ConvertText)

	log.Fatal(app.Listen(":3000"))

}
