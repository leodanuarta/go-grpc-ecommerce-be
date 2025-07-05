package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Post("/product/upload")
	app.Listen(":8081")
}
