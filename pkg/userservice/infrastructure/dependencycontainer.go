package infrastructure

import (
	"github.com/jmoiron/sqlx"
	"userservice/pkg/userservice/app/auth"
	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/domain"
	"userservice/pkg/userservice/infrastructure/mysql/repository"
)

type Parameters interface {
	HasherSalt() string
}

type DependencyContainer interface {
	UserService() service.UserService
	AuthenticationService() auth.AuthenticationService
}

func NewDependencyContainer(client *sqlx.DB, parameters Parameters) DependencyContainer {
	return &dependencyContainer{client: client, parameters: parameters}
}

type dependencyContainer struct {
	client     *sqlx.DB
	parameters Parameters
}

func (container *dependencyContainer) UserService() service.UserService {
	return service.NewUserService(
		container.DomainUserService(),
		container.Hasher(),
	)
}

func (container *dependencyContainer) AuthenticationService() auth.AuthenticationService {
	return auth.NewAuthenticationService(container.UserRepository(), container.Hasher())
}

func (container *dependencyContainer) DomainUserService() domain.UserService {
	return domain.NewUserService(container.UserRepository())
}

func (container *dependencyContainer) UserRepository() domain.UserRepository {
	return repository.NewUserRepository(container.client)
}

func (container *dependencyContainer) Hasher() hash.Hasher {
	return hash.NewSHA1Hasher(container.parameters.HasherSalt())
}
