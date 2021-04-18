package auth

import (
	"github.com/google/uuid"
	appservice "userservice/pkg/userservice/app/service"

	"github.com/pkg/errors"
	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/domain"
)

var ErrIncorrectAuthData = errors.New("incorrect auth data")

type AuthenticationService interface {
	AuthenticateUser(email, password string) (string, appservice.Role, error)
}

func NewAuthenticationService(userRepo domain.UserRepository, hasher hash.Hasher) AuthenticationService {
	return &authenticationService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

type authenticationService struct {
	domainService domain.UserService
	userRepo      domain.UserRepository
	hasher        hash.Hasher
}

func (service *authenticationService) AuthenticateUser(email, password string) (string, appservice.Role, error) {
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		return "", 0, err
	}

	if service.hasher.Hash(password) != user.Password {
		return "", 0, ErrIncorrectAuthData
	}

	return uuid.UUID(user.ID).String(), appservice.Role(user.Role), err
}
