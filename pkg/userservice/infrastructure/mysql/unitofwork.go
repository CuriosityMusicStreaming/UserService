package mysql

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/pkg/errors"

	"userservice/pkg/userservice/app/service"
	"userservice/pkg/userservice/domain"
	"userservice/pkg/userservice/infrastructure/mysql/repository"
)

func NewUnitOfFactory(client mysql.TransactionalClient) service.UnitOfWorkFactory {
	return &unitOfWorkFactory{client: client}
}

type unitOfWorkFactory struct {
	client mysql.TransactionalClient
}

func (factory *unitOfWorkFactory) NewUnitOfWork(_ string) (service.UnitOfWork, error) {
	transaction, err := factory.client.BeginTransaction()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &unitOfWork{transaction: transaction}, nil
}

type unitOfWork struct {
	transaction mysql.Transaction
}

func (u *unitOfWork) UserRepository() domain.UserRepository {
	return repository.NewUserRepository(u.transaction)
}

func (u *unitOfWork) Complete(err error) error {
	if err != nil {
		err2 := u.transaction.Rollback()
		if err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
	}

	return errors.WithStack(u.transaction.Commit())
}
