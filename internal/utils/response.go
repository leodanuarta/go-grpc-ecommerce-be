package utils

import (
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func ValidationErrorResponse(validationsError []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "Validation error",
		IsError:          true,
		ValidationErrors: validationsError,
	}
}

func BadRequestResponse(msg string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		IsError:    true,
		Message:    msg,
	}
}

func UnauthenticatedResponse() error {
	return status.Errorf(codes.Unauthenticated, "Unauthenticated")
}
