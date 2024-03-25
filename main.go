package main

import (
	httpHandler "auditlog/helper"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	http := httpHandler.InitHttp()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.All("/", http.HTTP)

	app.Use(cors.New())

	log.Fatal(app.Listen(":9090"))
}
