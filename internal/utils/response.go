package utils

import "github.com/leodanuarta/go-grpc-ecommerce-be/pb/common"

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
