package handler

import (
	"context"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/service"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
)

func (sh *authHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := sh.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := sh.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (sh *authHandler) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &auth.LogoutResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	// process register
	res, err := sh.authService.Logout(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type authHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
