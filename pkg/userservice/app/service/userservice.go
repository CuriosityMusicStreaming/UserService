package service

import (
	"github.com/google/uuid"

	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/domain"
)

type Role int

const (
	Listener = Role(domain.Listener)
	Creator  = Role(domain.Creator)
)

type UserService interface {
	AddUser(email, password string, role Role) (string, error)
}

func NewUserService(unitOfWorkFactory UnitOfWorkFactory, hasher hash.Hasher) UserService {
	return &userService{
		unitOfWorkFactory: unitOfWorkFactory,
		hasher:            hasher,
	}
}

type userService struct {
	unitOfWorkFactory UnitOfWorkFactory
	hasher            hash.Hasher
}

func (service *userService) AddUser(email, password string, role Role) (string, error) {
	var userID domain.UserID

	err := service.executeInUnitOfWork(func(provider RepositoryProvider) error {
		domainService := domain.NewUserService(provider.UserRepository())

		var err2 error

		userID, err2 = domainService.AddUser(email, service.hasher.Hash(password), domain.Role(role))

		return err2
	})

	if err != nil {
		return "", err
	}

	return uuid.UUID(userID).String(), nil
}

func (service *userService) executeInUnitOfWork(f func(provider RepositoryProvider) error) error {
	unitOfWork, err := service.unitOfWorkFactory.NewUnitOfWork("")
	if err != nil {
		return err
	}
	defer func() {
		err = unitOfWork.Complete(err)
	}()
	err = f(unitOfWork)
	return err
}
