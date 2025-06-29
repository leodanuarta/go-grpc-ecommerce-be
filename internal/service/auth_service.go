package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/repository"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (as *authService) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if request.Password != request.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Passowrd is not matched"),
		}, nil
	}

	// ngecek email ke database
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	// apabila email sudah terdaftar, kita error in
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exists"),
		}, nil
	}
	// HASH password nya
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, err
	}

	// jika belom, insert ke DB
	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  string(hasedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &request.FullName,
	}

	err = as.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User is registered"),
	}, nil
}

// Login implements IAuthService.
func (as *authService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	// check apakah email ada
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}

	// check apakah password sama
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
		}
	}

	// generate jwt
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})

	secretKey := os.Getenv("JWT_SECRET")

	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	// kirim response
	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login successful"),
		AccessToken: accessToken,
	}, nil
}

// Logout implements IAuthService.
func (as *authService) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// dapatkan token dari metadata grpc
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	bearerToken, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	if len(bearerToken) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	tokenSplit := strings.Split(bearerToken[0], " ")

	if len(tokenSplit) != 2 {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	if tokenSplit[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	jwtToken := tokenSplit[1]

	// kembalikan token tadi hingga menjadi entity jwt
	tokenClaims, err := jwt.ParseWithClaims(jwtToken, &entity.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if !tokenClaims.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	var claims *entity.JwtClaims

	if claims, ok = tokenClaims.Claims.(*entity.JwtClaims); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	// kita masukkan token dari metadata ke dalam memori db / cache
	as.cacheService.Set(jwtToken, "", time.Duration(claims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)
	// kirim response

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout Successful"),
	}, nil
}

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error)
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
