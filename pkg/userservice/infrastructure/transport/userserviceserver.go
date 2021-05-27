package transport

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	api "userservice/api/userservice"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/infrastructure"
)

func NewUserServiceServer(container infrastructure.DependencyContainer) api.UserServiceServer {
	return &userServiceServer{container: container}
}

type userServiceServer struct {
	container infrastructure.DependencyContainer
}

func (server *userServiceServer) AddUser(_ context.Context, req *api.AddUserRequest) (*api.AddUserResponse, error) {
	role, ok := apiToUserRoleMap[req.Role]
	if !ok {
		return nil, ErrUnknownUserRole
	}

	userID, err := server.container.UserService().AddUser(req.Email, req.Password, role)
	if err != nil {
		return nil, err
	}

	return &api.AddUserResponse{UserId: userID}, nil
}

var apiToUserRoleMap = map[api.UserRole]service.Role{
	api.UserRole_LISTENER: service.Listener,
	api.UserRole_CREATOR:  service.Creator,
}

var userRoleToApiMap = map[service.Role]api.UserRole{
	service.Listener: api.UserRole_LISTENER,
	service.Creator:  api.UserRole_CREATOR,
}

var (
	ErrUnknownUserRole = errors.New("unknown user role")
)
