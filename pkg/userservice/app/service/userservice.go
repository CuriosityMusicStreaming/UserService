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

func NewUserService(domainService domain.UserService, hasher hash.Hasher) UserService {
	return &userService{
		domainService: domainService,
		hasher:        hasher,
	}
}

type userService struct {
	domainService domain.UserService
	hasher        hash.Hasher
}

func (service *userService) AddUser(email, password string, role Role) (string, error) {
	userId, err := service.domainService.AddUser(email, service.hasher.Hash(password), domain.Role(role))
	if err != nil {
		return "", err
	}

	return uuid.UUID(userId).String(), nil
}
