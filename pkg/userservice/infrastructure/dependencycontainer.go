package infrastructure

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"userservice/pkg/userservice/app/auth"
	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/domain"
	mysql2 "userservice/pkg/userservice/infrastructure/mysql"
	"userservice/pkg/userservice/infrastructure/mysql/repository"
)

type Parameters interface {
	HasherSalt() string
}

type DependencyContainer interface {
	UserService() service.UserService
	AuthenticationService() auth.AuthenticationService
}

func NewDependencyContainer(client mysql.TransactionalClient, parameters Parameters) DependencyContainer {
	return &dependencyContainer{
		client:     client,
		parameters: parameters,
		unitOfWork: unitOfWorkFactory(client),
	}
}

type dependencyContainer struct {
	client     mysql.TransactionalClient
	unitOfWork service.UnitOfWorkFactory
	parameters Parameters
}

func (container *dependencyContainer) UserService() service.UserService {
	return service.NewUserService(
		container.unitOfWork,
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

func unitOfWorkFactory(client mysql.TransactionalClient) service.UnitOfWorkFactory {
	return mysql2.NewUnitOfFactory(client)
}
