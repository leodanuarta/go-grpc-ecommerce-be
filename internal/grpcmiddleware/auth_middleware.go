package grpcmiddleware

import (
	"context"
	"log"

	jwtEntity "github.com/leodanuarta/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"google.golang.org/grpc"

	gocache "github.com/patrickmn/go-cache"
)

var publicApis = map[string]bool{
	"/auth.AuthService/Login":               true,
	"/auth.AuthService/Register":            true,
	"/product.ProductService/DetailProduct": true,
	"/product.ProductService/ListProduct":   true,
}

func (am *authMiddleware) AuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	log.Println(info.FullMethod)

	if publicApis[info.FullMethod] {
		return handler(ctx, req)
	}

	// Ambil token dari metadata
	tokenStr, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// cek token dari logout cache
	_, ok := am.cacheService.Get(tokenStr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// Parse jwt nya hingga menjadi entity
	claims, err := jwtEntity.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, err
	}

	// sematkan entity ke context
	ctx = claims.SetToContext(ctx)

	response, err := handler(ctx, req)

	return response, err
}

type authMiddleware struct {
	cacheService *gocache.Cache
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
