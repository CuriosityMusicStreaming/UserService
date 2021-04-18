package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Role int

const (
	Listener Role = iota
	Creator
)

type UserID uuid.UUID

type User struct {
	ID       UserID
	Email    string
	Password string
	Role
}

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrUserWithEmailAlreadyExists = errors.New("user with email already exists")
)

type UserRepository interface {
	NewID() UserID
	Find(id UserID) (User, error)
	FindByEmail(email string) (User, error)
	Store(user User) error
	Remove(id UserID) error
}
