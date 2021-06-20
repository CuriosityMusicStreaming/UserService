package transport

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"userservice/pkg/userservice/app/auth"
	"userservice/pkg/userservice/app/query"
	"userservice/pkg/userservice/domain"
)

func translateError(err error) error {
	switch errors.Cause(err) {
	case auth.ErrOnlyCreatorsCanAddContent:
		return status.Error(codes.PermissionDenied, err.Error())
	case domain.ErrUserNotFound:
	case query.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrUserWithEmailAlreadyExists:
	case auth.ErrIncorrectAuthData:
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return err
}
