package main

import (
	httpHandler "auditlog/helper"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/utils"
	"time"
)

func main() {
	http := httpHandler.InitHttp()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(
		helmet.New(),
	)

	csrfConfig := csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "dp_csrf_id",
		CookieSameSite: "Strict",
		Expiration:     3 * time.Hour,
		KeyGenerator:   utils.UUID,
	}

	app.Use(
		csrf.New(csrfConfig),
	)

	limiterConfig := limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "0.0.0.0"
		},
		Max:        20,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "too fast")
		},
	}

	app.Use(
		limiter.New(limiterConfig),
	)

	//app.Use(
	//	logger.New(),
	//)

	app.Use(cors.New())

	app.All("/", http.HTTP)
	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":9090"))
}
