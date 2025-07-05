package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/common"
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

func (p *productRepository) GetProductPagination(ctx context.Context, pagination *common.PaginationRequest) ([]*entity.Product, *common.PaginationResponse, error) {
	offset := (pagination.CurrentPage - 1) * pagination.ItemPerPage

	row := p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM product WHERE is_deleted = false")
	if row.Err() != nil {
		return nil, nil, row.Err()
	}

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalCount + int(pagination.ItemPerPage) - 1) / int(pagination.ItemPerPage)

	rows, err := p.db.QueryContext(
		ctx,
		`SELECT id, name, description, price, image_file_name
			FROM product
			WHERE is_deleted = false
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
			`,
		pagination.ItemPerPage,
		offset,
	)

	if err != nil {
		return nil, nil, err
	}

	var products []*entity.Product = make([]*entity.Product, 0)
	for rows.Next() {
		var product entity.Product

		err := rows.Scan(&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageFileName,
		)

		if err != nil {
			return nil, nil, err
		}

		products = append(products, &product)
	}

	paginationResp := &common.PaginationResponse{
		CurrentPage:    pagination.CurrentPage,
		ItemPerPage:    pagination.ItemPerPage,
		TotalItemCount: int32(totalCount),
		TotalPageCount: int32(totalPages),
	}
	return products, paginationResp, err
}
