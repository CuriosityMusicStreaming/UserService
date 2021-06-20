package auth

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/pkg/errors"

	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/app/query"
	appservice "userservice/pkg/userservice/app/service"
)

var (
	ErrIncorrectAuthData         = errors.New("incorrect auth data")
	ErrOnlyCreatorsCanAddContent = errors.New("only creators can add content")
)

type AuthenticationService interface {
	AuthenticateUser(email, password string) (string, appservice.Role, error)
	CanAddContent(descriptor auth.UserDescriptor) (bool, error)
}

func NewAuthenticationService(queryService query.UserQueryService, hasher hash.Hasher) AuthenticationService {
	return &authenticationService{
		queryService: queryService,
		hasher:       hasher,
	}
}

type authenticationService struct {
	queryService query.UserQueryService
	hasher       hash.Hasher
}

func (service *authenticationService) AuthenticateUser(email, password string) (string, appservice.Role, error) {
	user, err := service.queryService.GetByEmail(email)
	if err != nil {
		return "", 0, err
	}

	if service.hasher.Hash(password) != user.Password {
		return "", 0, ErrIncorrectAuthData
	}

	return user.ID.String(), appservice.Role(user.Role), err
}

func (service *authenticationService) CanAddContent(userDescriptor auth.UserDescriptor) (bool, error) {
	user, err := service.queryService.GetUser(userDescriptor.UserID)
	if err != nil {
		return false, err
	}

	canAdd := user.Role == query.Creator
	if !canAdd {
		return false, ErrOnlyCreatorsCanAddContent
	}

	return true, nil
}
