package handler

import (
	"context"
	"fmt"

	"github.com/leodanuarta/go-grpc-ecommerce-be/internal/utils"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/service"
)

type IServiceHandler interface {
	HelloWorld(ctx context.Context, request *service.HelloWorldRequest) (*service.HelloWorldResponse, error)
}

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, request *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	validationError, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}

	if validationError != nil {
		return &service.HelloWorldResponse{
			Base: utils.ValidationErrorResponse(validationError),
		}, nil
	}

	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s ", request.Name),
		Base:    utils.SuccessResponse("Success"),
	}, nil
}
func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}
