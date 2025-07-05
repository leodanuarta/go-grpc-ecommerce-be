package handler

import (
	service "github.com/leodanuarta/go-grpc-ecommerce-be/internal/service/product"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/product"
)

type productHandler struct {
	product.UnimplementedProductServiceServer

	productService service.IProductService
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
}
