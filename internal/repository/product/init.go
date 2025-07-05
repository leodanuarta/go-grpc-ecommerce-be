package repository

import (
	"context"
	"database/sql"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
)

type IProductRepository interface {
	CreateNewProduct(ctx context.Context, product *entity.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}
