package jwt

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type JwtEntityContextKey string

var JwtEntityContextKeyValue JwtEntityContextKey = "JwtEntity"

func GetClaimsFromToken(token string) (*JwtClaims, error) {
	// if token == "" {
	// 	return nil, utils.UnauthenticatedResponse()
	// }
	tokenClaims, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	if !tokenClaims.Valid {
		return nil, utils.UnauthenticatedResponse()
	}

	if claims, ok := tokenClaims.Claims.(*JwtClaims); ok {
		return claims, nil
	}

	return nil, utils.UnauthenticatedResponse()
}

func (jc *JwtClaims) SetToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, JwtEntityContextKeyValue, jc)

}

func GetClaimsFromContext(ctx context.Context) (*JwtClaims, error) {
	token, err := ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := GetClaimsFromToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
