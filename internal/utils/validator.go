package utils

import (
	"errors"

	"buf.build/go/protovalidate"
	"github.com/leodanuarta/go-grpc-ecommerce-be/pb/common"
	"google.golang.org/protobuf/proto"
)

func CheckValidation(req proto.Message) ([]*common.ValidationError, error) {
	if err := protovalidate.Validate(req); err != nil {
		var validationError *protovalidate.ValidationError
		if errors.As(err, &validationError) {
			var validationErrorResponse []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, v := range validationError.Violations {
				validationErrorResponse = append(validationErrorResponse, &common.ValidationError{
					Message: *v.Proto.Message,
				})
			}
			return validationErrorResponse, nil
		}

		return nil, err
	}

	return nil, nil
}
