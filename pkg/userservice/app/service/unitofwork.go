package service

import "userservice/pkg/userservice/domain"

type UnitOfWorkFactory interface {
	NewUnitOfWork(lockName string) (UnitOfWork, error)
}

type RepositoryProvider interface {
	UserRepository() domain.UserRepository
}

type UnitOfWork interface {
	RepositoryProvider
	Complete(err error) error
}
