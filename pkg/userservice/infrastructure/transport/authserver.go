package transport

import (
	"golang.org/x/net/context"
	authenticationapi "userservice/api/authenticationservice"
	authorizationapi "userservice/api/authorizationservice"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/infrastructure"
)

func NewAuthServer(container infrastructure.DependencyContainer) *authServer {
	return &authServer{container: container}
}

type authServer struct {
	container infrastructure.DependencyContainer
}

func (server *authServer) CanAddContent(_ context.Context, req *authorizationapi.CanAddContentRequest) (*authorizationapi.CanAddContentResponse, error) {
	userDesc, err := server.container.UserDescriptorSerializer().Deserialize(req.UserToken)
	if err != nil {
		return nil, err
	}

	canAddContent, err := server.container.AuthenticationService().CanAddContent(userDesc)
	if err != nil {
		return nil, err
	}

	return &authorizationapi.CanAddContentResponse{CanAdd: canAddContent}, nil
}

func (server *authServer) AuthenticateUser(_ context.Context, req *authenticationapi.AuthenticateUserRequest) (*authenticationapi.AuthenticateUserResponse, error) {
	userID, role, err := server.container.AuthenticationService().AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &authenticationapi.AuthenticateUserResponse{
		UserID: userID,
		Role:   userRoleToAuthApiMap[role],
	}, nil
}

var userRoleToAuthApiMap = map[service.Role]authenticationapi.UserRole{
	service.Listener: authenticationapi.UserRole_LISTENER,
	service.Creator:  authenticationapi.UserRole_CREATOR,
}
