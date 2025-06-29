package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/handler"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/repository"
	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/service"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/auth"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pkg/database"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pkg/grpcmiddleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	ctx := context.Background()

	godotenv.Load()

	list, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Panicf("Error when listening %v", err)
	}

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to Database ...")

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection is registered.")
	}

	log.Println("Server is running on port : 50052 port")

	if err := serv.Serve(list); err != nil {
		log.Panicf("server is error : %v", err)
	}
}
