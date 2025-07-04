package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	jwtEntity "github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity/jwt"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtEntity.JwtClaims{
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
	jwtToken, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// kembalikan token tadi hingga menjadi entity jwt
	tokenClaims, err := jwtEntity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	// kita masukkan token dari metadata ke dalam memori db / cache
	as.cacheService.Set(jwtToken, "", time.Duration(tokenClaims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)
	// kirim response

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout Successful"),
	}, nil
}

// ChangePassword implements IAuthService.
func (as *authService) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	// Cek apakah new passwd confirmmation matched ?
	if request.NewPassword != request.NewPasswordConfirmation {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("New password is not matched"),
		}, nil
	}

	// cek apakah old password matched ?
	jwtToken, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := jwtEntity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("User doesn't exists"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &auth.ChangePasswordResponse{
				Base: utils.BadRequestResponse("old password is not matched"),
			}, nil
		}
	}

	// Update new password ke database
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10)
	if err != nil {
		return nil, err
	}

	err = as.authRepository.UpdateUserPassword(ctx, user.Id, string(hashedNewPassword), user.FullName)
	if err != nil {
		return nil, err
	}

	// kirim response
	return &auth.ChangePasswordResponse{
		Base: utils.SuccessResponse("Chnage passoword success"),
	}, nil

}

// GetProfile implements IAuthService.
func (as *authService) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	// Get data token
	claims, ok := ctx.Value(jwtEntity.JwtEntityContextKeyValue).(*jwtEntity.JwtClaims)
	if !ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// Ambil data dari database
	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.GetProfileResponse{
			Base: utils.BadRequestResponse("user doesn't exists"),
		}, nil
	}

	// Buat response

	return &auth.GetProfileResponse{
		Base:        utils.SuccessResponse("Get profile success"),
		UserId:      claims.Subject,
		FullName:    claims.FullName,
		Email:       claims.Email,
		RoleCode:    claims.Role,
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil

}
