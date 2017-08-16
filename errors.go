package alecton

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewNotFoundError(kind, name string) error {
	return grpc.Errorf(codes.NotFound, "not found: %s/%s", kind, name)
}

func IsNotFoundError(err error) bool {
	return isErrorCode(err, codes.NotFound)
}

func NewAlreadyExistsError(kind, name string) error {
	return grpc.Errorf(codes.AlreadyExists, "already exists: %s/%s", kind, name)
}

func IsAlreadyExistsError(err error) bool {
	return isErrorCode(err, codes.AlreadyExists)
}

func isErrorCode(err error, c codes.Code) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return s.Code() == c
}
