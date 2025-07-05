package handler

import (
	service "github.com/leodanuarta/go-grpc-ecommerce-be/internal/service/auth"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
