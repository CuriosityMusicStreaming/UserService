package infrastructure

import (
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	commonmysql "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"userservice/pkg/userservice/app/auth"
	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/app/query"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/domain"
	"userservice/pkg/userservice/infrastructure/mysql"
	mysqlquery "userservice/pkg/userservice/infrastructure/mysql/query"
	"userservice/pkg/userservice/infrastructure/mysql/repository"
)

type Parameters interface {
	HasherSalt() string
}

type DependencyContainer interface {
	UserService() service.UserService
	AuthenticationService() auth.AuthenticationService
	UserDescriptorSerializer() commonauth.UserDescriptorSerializer
}

func NewDependencyContainer(client commonmysql.TransactionalClient, parameters Parameters) DependencyContainer {
	return &dependencyContainer{
		client:     client,
		parameters: parameters,
		unitOfWork: unitOfWorkFactory(client),
	}
}

type dependencyContainer struct {
	client     commonmysql.TransactionalClient
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
	return auth.NewAuthenticationService(container.userQueryService(), container.Hasher())
}

func (container *dependencyContainer) UserDescriptorSerializer() commonauth.UserDescriptorSerializer {
	return commonauth.NewUserDescriptorSerializer()
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

func (container *dependencyContainer) userQueryService() query.UserQueryService {
	return mysqlquery.NewUserQueryService(container.client)
}

func unitOfWorkFactory(client commonmysql.TransactionalClient) service.UnitOfWorkFactory {
	return mysql.NewUnitOfFactory(client)
}
