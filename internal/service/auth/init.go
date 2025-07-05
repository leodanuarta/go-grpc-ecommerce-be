package service

import (
	"context"

	repository "github.com/leodanuarta/go-grpc-ecommerce-be/internal/repository/auth"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
	gocache "github.com/patrickmn/go-cache"
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache
}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}
