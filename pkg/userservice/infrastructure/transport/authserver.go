package transport

import (
	"golang.org/x/net/context"

	authenticationapi "userservice/api/authenticationservice"
	authorizationapi "userservice/api/authorizationservice"
	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/infrastructure"
)

type AuthServer interface {
	authorizationapi.AuthorizationServiceServer
	authenticationapi.AuthenticationServiceServer
}

func NewAuthServer(container infrastructure.DependencyContainer) AuthServer {
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
		Role:   userRoleToAuthAPIMap[role],
	}, nil
}

var userRoleToAuthAPIMap = map[service.Role]authenticationapi.UserRole{
	service.Listener: authenticationapi.UserRole_LISTENER,
	service.Creator:  authenticationapi.UserRole_CREATOR,
}
