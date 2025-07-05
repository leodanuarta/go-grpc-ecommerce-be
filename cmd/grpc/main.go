package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/grpcmiddleware"
	handlerAuth "github.com/leodanuarta/go-grpc-ecommerce-be/internal/handler/auth"
	handlerProduct "github.com/leodanuarta/go-grpc-ecommerce-be/internal/handler/product"
	repositoryAuth "github.com/leodanuarta/go-grpc-ecommerce-be/internal/repository/auth"
	repositoryProduct "github.com/leodanuarta/go-grpc-ecommerce-be/internal/repository/product"
	serviceAuth "github.com/leodanuarta/go-grpc-ecommerce-be/internal/service/auth"
	serviceProduct "github.com/leodanuarta/go-grpc-ecommerce-be/internal/service/product"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/product"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pkg/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	gocache "github.com/patrickmn/go-cache"
)

func main() {

	ctx := context.Background()

	godotenv.Load()

	list, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Panicf("Error when listening %v", err)
	}

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to Database ...")

	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware := grpcmiddleware.NewAuthMiddleware(cacheService)

	authRepository := repositoryAuth.NewAuthRepository(db)
	authService := serviceAuth.NewAuthService(authRepository, cacheService)
	authHandler := handlerAuth.NewAuthHandler(authService)

	productRepository := repositoryProduct.NewProductRepository(db)
	productService := serviceProduct.NewProductService(productRepository)
	productHandler := handlerProduct.NewProductHandler(productService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
			authMiddleware.AuthMiddleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductServiceServer(serv, productHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection is registered.")
	}

	log.Println("Server is running on port : 50053 port")

	if err := serv.Serve(list); err != nil {
		log.Panicf("server is error : %v", err)
	}
}
