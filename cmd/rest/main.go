package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	handler "github.com/leodanuarta/go-grpc-ecommerce-be/internal/handler/product"
)

func handleGetFileName(c *fiber.Ctx) error {
	fileNameParam := c.Params("filename")
	filePath := path.Join("storage", "product", fileNameParam)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			log.Println(err)
			return c.Status(http.StatusNotFound).SendString("Not Found")
		}

		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// buka file nya
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	c.Set("Content-Type", mimeType)
	return c.SendStream(file)
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/storage/product/:filename", handleGetFileName)
	app.Post("/product/upload", handler.UploadProductImageHandler)
	app.Listen(":8081")
}
