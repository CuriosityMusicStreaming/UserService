package transport

import (
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	api "userservice/api/userservice"
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

func (server *userServiceServer) AuthenticateUser(_ context.Context, req *api.AuthenticateUserRequest) (*api.AuthenticateUserResponse, error) {
	userID, role, err := server.container.AuthenticationService().AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &api.AuthenticateUserResponse{
		UserID: userID,
		Role:   userRoleToApiMap[role],
	}, nil
}

var apiToUserRoleMap = map[api.UserRole]commonauth.Role{
	api.UserRole_LISTENER: commonauth.Listener,
	api.UserRole_CREATOR:  commonauth.Creator,
}

var userRoleToApiMap = map[commonauth.Role]api.UserRole{
	commonauth.Listener: api.UserRole_LISTENER,
	commonauth.Creator:  api.UserRole_CREATOR,
}

var (
	ErrUnknownUserRole = errors.New("unknown user role")
)
