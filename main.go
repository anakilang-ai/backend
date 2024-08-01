package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/anakilang-ai/backend/routes"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return routes.URL(c)
	})
	port := ":8080"
	fmt.Println("Server started at: http://localhost" + port)
	app.Listen(port)
}