package grpcmiddleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack()

			err = status.Errorf(codes.Internal, "Internal Server Error")
		}
	}()
	response, err := handler(ctx, req)
	if err != nil {
		log.Println(err)

		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.Unauthenticated {
				return nil, err
			}
		}

		// return nil, status.Errorf(codes.Internal, "Internal Server Error")
		return nil, err
	}

	return response, err
}
