package infrastructure

import (
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	commonmysql "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"

	"userservice/pkg/userservice/app/auth"
	"userservice/pkg/userservice/app/hash"
	"userservice/pkg/userservice/app/query"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/infrastructure/mysql"
	mysqlquery "userservice/pkg/userservice/infrastructure/mysql/query"
)

type Parameters interface {
	HasherSalt() string
}

type DependencyContainer interface {
	UserService() service.UserService
	AuthenticationService() auth.AuthenticationService
	UserDescriptorSerializer() commonauth.UserDescriptorSerializer
	UserQueryService() query.UserQueryService
}

func NewDependencyContainer(client commonmysql.TransactionalClient, parameters Parameters) DependencyContainer {
	userQueryService := userQueryService(client)
	hasher := hasher(parameters)

	return &dependencyContainer{
		userService:              userService(unitOfWorkFactory(client), hasher),
		authenticationService:    authenticationService(userQueryService, hasher),
		userDescriptorSerializer: userDescriptorSerializer(),
	}
}

type dependencyContainer struct {
	userService              service.UserService
	userQueryService         query.UserQueryService
	authenticationService    auth.AuthenticationService
	userDescriptorSerializer commonauth.UserDescriptorSerializer
}

func (container *dependencyContainer) UserService() service.UserService {
	return container.userService
}

func (container *dependencyContainer) AuthenticationService() auth.AuthenticationService {
	return container.authenticationService
}

func (container *dependencyContainer) UserDescriptorSerializer() commonauth.UserDescriptorSerializer {
	return container.userDescriptorSerializer
}

func (container *dependencyContainer) UserQueryService() query.UserQueryService {
	return container.userQueryService
}

func userService(unitOfWorkFactory service.UnitOfWorkFactory, hasher hash.Hasher) service.UserService {
	return service.NewUserService(
		unitOfWorkFactory,
		hasher,
	)
}

func authenticationService(queryService query.UserQueryService, hasher hash.Hasher) auth.AuthenticationService {
	return auth.NewAuthenticationService(queryService, hasher)
}

func userDescriptorSerializer() commonauth.UserDescriptorSerializer {
	return commonauth.NewUserDescriptorSerializer()
}

func hasher(parameters Parameters) hash.Hasher {
	return hash.NewSHA1Hasher(parameters.HasherSalt())
}

func userQueryService(client commonmysql.TransactionalClient) query.UserQueryService {
	return mysqlquery.NewUserQueryService(client)
}

func unitOfWorkFactory(client commonmysql.TransactionalClient) service.UnitOfWorkFactory {
	return mysql.NewUnitOfFactory(client)
}
