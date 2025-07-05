package service

import (
	"context"

	repository "github.com/leodanuarta/go-grpc-ecommerce-be/internal/repository/product"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/product"
)

type IProductService interface {
	CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
}

func NewProductService(productRepository repository.IProductRepository) IProductService {
	return &productService{
		productRepository: productRepository,
	}
}
