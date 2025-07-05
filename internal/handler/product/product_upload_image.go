package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UploadProductImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image data not found",
		})
	}

	// validasi gambar
	// validasi extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "image extension is not allowed (jpg, jpeg, png. webp)",
		})
	}

	// validasi content type
	contentType := file.Header.Get("Content-Type")
	allowedContentType := map[string]bool{
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedContentType[contentType] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "content type is not allowed",
		})
	}

	// product_160722.png
	timestampt := time.Now().UnixNano()
	filename := fmt.Sprintf("product_%d%s", timestampt, filepath.Ext(file.Filename))
	uploadPath := "./storage/product/" + filename
	err = c.SaveFile(file, uploadPath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Upload Success",
		"file_name": filename,
	})
}
