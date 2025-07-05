package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/product"

	jwtEntity "github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity/jwt"
)

func (ps *productService) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	// cek dulu apakah user nya merupakan admin
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}
	// apakah image nya ada ?
	imagePath := filepath.Join("storage", "product", request.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &product.CreateProductResponse{
				Base: utils.BadRequestResponse("file not found"),
			}, nil
		}
	}

	// insert ke database
	productEntity := entity.Product{
		Id:            uuid.NewString(),
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		ImageFileName: request.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     claims.FullName,
	}
	err = ps.productRepository.CreateNewProduct(ctx, &productEntity)
	if err != nil {
		return nil, err
	}

	// success
	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Product is created"),
		Id:   productEntity.Id,
	}, nil
}

// DetailProduct implements IProductService.
func (ps *productService) DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	// query ke db dengan data ID
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	// apabila null, kita return not found
	if productEntity == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not Found"),
		}, nil
	}

	// kirim respon detail
	return &product.DetailProductResponse{
		Base:        utils.SuccessResponse("Get product detail success"),
		Id:          productEntity.Id,
		Name:        productEntity.Name,
		Price:       productEntity.Price,
		Description: productEntity.Description,
		ImageUrl:    fmt.Sprintf("%s/product/%s", os.Getenv("STORAGE_SERVICE_URL"), productEntity.ImageFileName),
	}, nil
}

func (ps *productService) EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error) {
	// cek dulu apakah user nya merupakan admin
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}

	// validasi apakah id yang dikirim itu ada di db ?
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.EditProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}
	// kalau gambarnya berubah, hapus gambar lama
	if productEntity.ImageFileName != request.ImageFileName {
		newImagePath := filepath.Join("storage", "product", request.ImageFileName)
		_, err = os.Stat(newImagePath)
		if err != nil {
			if os.IsNotExist(err) {
				return &product.EditProductResponse{
					Base: utils.BadRequestResponse("Image not found"),
				}, nil
			}

			return nil, err
		}

		oldImagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
		err = os.Remove(oldImagePath)
		if err != nil {
			return nil, err
		}
	}

	// update ke database
	tempProduct := entity.Product{
		Id:            request.Id,
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		ImageFileName: request.ImageFileName,
		UpdatedAt:     time.Now(),
		UpdatedBy:     &claims.FullName,
	}
	err = ps.productRepository.UpdateProduct(ctx, &tempProduct)
	if err != nil {
		return nil, err
	}

	// kirim respon detail
	return &product.EditProductResponse{
		Base: utils.SuccessResponse("Edit product success"),
		Id:   request.Id,
	}, nil
}

// DeleteProduct implements IProductService.
func (ps *productService) DeleteProduct(ctx context.Context, request *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	// cek dulu apakah user nya merupakan admin
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if claims.Role != entity.UserRoleAdmin {
		return nil, utils.UnauthenticatedResponse()
	}
	// validasi apakah id yang dikirim itu ada di db ?
	productEntity, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if productEntity == nil {
		return &product.DeleteProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	err = ps.productRepository.DeleteProduct(ctx, request.Id, time.Now(), claims.FullName)
	if err != nil {
		return nil, err
	}

	imagePath := filepath.Join("storage", "product", productEntity.ImageFileName)
	err = os.Remove(imagePath)
	if err != nil {
		return nil, err
	}

	return &product.DeleteProductResponse{
		Base: utils.SuccessResponse("Delete product success"),
	}, nil
}
