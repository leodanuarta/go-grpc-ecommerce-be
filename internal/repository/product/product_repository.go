package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
)

// CreateNewProduct implements IProductRepository.
func (p *productRepository) CreateNewProduct(ctx context.Context, product *entity.Product) error {
	_, err := p.db.ExecContext(
		ctx,
		`INSERT INTO product (id, name, description, price, image_file_name, created_at, created_by, updated_at, 
				updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		product.Id,
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.CreatedAt,
		product.CreatedBy,
		product.UpdatedAt,
		product.UpdatedBy,
		product.DeletedAt,
		product.DeletedBy,
		product.IsDeleted,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) GetProductById(ctx context.Context, id string) (*entity.Product, error) {
	var productEntty entity.Product
	row := p.db.QueryRowContext(
		ctx,
		"SELECT id, name, description, price, image_file_name FROM product WHERE id = $1 AND is_deleted = false",
		id,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&productEntty.Id,
		&productEntty.Name,
		&productEntty.Description,
		&productEntty.Price,
		&productEntty.ImageFileName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &productEntty, nil
}

// EditProduct implements IProductRepository.
func (p *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	_, err := p.db.ExecContext(
		ctx,
		`UPDATE product SET
			name=$1,
			description=$2,
			price=$3,
			image_file_name=$4,
			updated_at=$5,
			updated_by=$6
			WHERE id = $7
			`,
		product.Name,
		product.Description,
		product.Price,
		product.ImageFileName,
		product.UpdatedAt,
		product.UpdatedBy,
		product.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *productRepository) DeleteProduct(ctx context.Context, id string, deletedAt time.Time, deleteBy string) error {
	_, err := p.db.ExecContext(
		ctx,
		`UPDATE product SET
			deleted_at = $1,
			deleted_by = $2,
			is_deleted = true
			WHERE id = $3
			`,
		deletedAt,
		deleteBy,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
