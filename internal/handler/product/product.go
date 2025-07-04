package handler

import (
	"context"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/product"
)

func (ph *productHandler) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := ph.productService.CreateProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &product.DetailProductResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := ph.productService.DetailProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) EditProduct(ctx context.Context, request *product.EditProductRequest) (*product.EditProductResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &product.EditProductResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := ph.productService.EditProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) DeleteProduct(ctx context.Context, request *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &product.DeleteProductResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := ph.productService.DeleteProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ph *productHandler) ListProduct(ctx context.Context, request *product.ListProductRequest) (*product.ListProductResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &product.ListProductResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := ph.productService.ListProduct(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
