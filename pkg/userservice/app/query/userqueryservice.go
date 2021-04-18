package query

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"userservice/pkg/userservice/domain"
)

type Role int

const (
	Listener = Role(domain.Listener)
	Creator  = Role(domain.Creator)
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserView struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     Role
}

type UserQueryService interface {
	GetUser(id uuid.UUID) (UserView, error)
	GetByEmail(email string) (UserView, error)
}
