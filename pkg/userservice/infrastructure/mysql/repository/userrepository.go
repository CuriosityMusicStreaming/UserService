package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"userservice/pkg/userservice/domain"
)

func NewUserRepository(client *sqlx.DB) domain.UserRepository {
	return &userRepository{client: client}
}

type userRepository struct {
	client *sqlx.DB
}

func (repo *userRepository) NewID() domain.UserID {
	return domain.UserID(uuid.New())
}

func (repo *userRepository) Find(id domain.UserID) (domain.User, error) {
	const selectSql = `SELECT * from user WHERE user_id = ?`

	binaryUUID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return domain.User{}, err
	}

	var user sqlxUser

	err = repo.client.Get(&user, selectSql, binaryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, errors.WithStack(err)
	}

	return domain.User{
		Id:       domain.UserID(user.UserId),
		Email:    user.Email,
		Password: user.Password,
		Role:     domain.Role(user.Role),
	}, nil
}

func (repo *userRepository) FindByEmail(email string) (domain.User, error) {
	const selectSql = `SELECT * from user WHERE email = ?`

	var user sqlxUser

	err := repo.client.Get(&user, selectSql, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, errors.WithStack(err)
	}

	return domain.User{
		Id:       domain.UserID(user.UserId),
		Email:    user.Email,
		Password: user.Password,
		Role:     domain.Role(user.Role),
	}, nil
}

func (repo *userRepository) Store(user domain.User) error {
	const insertSql = `INSERT INTO user VALUES(?, ?, ?, ?)`

	binaryUUID, err := uuid.UUID(user.Id).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(insertSql, binaryUUID, user.Email, user.Password, int(user.Role))
	return err
}

func (repo *userRepository) Remove(id domain.UserID) error {
	const deleteSql = `DELETE FROM user WHERE user_id = ?`

	binaryUUID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(deleteSql, binaryUUID)
	return err
}

type sqlxUser struct {
	UserId   uuid.UUID `db:"user_id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     int       `db:"role"`
}
