package repository

import (
	"database/sql"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"userservice/pkg/userservice/domain"
)

func NewUserRepository(client mysql.Client) domain.UserRepository {
	return &userRepository{client: client}
}

type userRepository struct {
	client mysql.Client
}

func (repo *userRepository) NewID() domain.UserID {
	return domain.UserID(uuid.New())
}

func (repo *userRepository) Find(id domain.UserID) (domain.User, error) {
	const selectSQL = `SELECT * from user WHERE user_id = ?`

	binaryUUID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return domain.User{}, err
	}

	var user sqlxUser

	err = repo.client.Get(&user, selectSQL, binaryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, errors.WithStack(err)
	}

	return domain.User{
		ID:       domain.UserID(user.UserID),
		Email:    user.Email,
		Password: user.Password,
		Role:     domain.Role(user.Role),
	}, nil
}

func (repo *userRepository) FindByEmail(email string) (domain.User, error) {
	const selectSQL = `SELECT * from user WHERE email = ?`

	var user sqlxUser

	err := repo.client.Get(&user, selectSQL, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, errors.WithStack(err)
	}

	return domain.User{
		ID:       domain.UserID(user.UserID),
		Email:    user.Email,
		Password: user.Password,
		Role:     domain.Role(user.Role),
	}, nil
}

func (repo *userRepository) Store(user domain.User) error {
	const insertSQL = `INSERT INTO user VALUES(?, ?, ?, ?)`

	binaryUUID, err := uuid.UUID(user.ID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(insertSQL, binaryUUID, user.Email, user.Password, int(user.Role))
	return err
}

func (repo *userRepository) Remove(id domain.UserID) error {
	const deleteSQL = `DELETE FROM user WHERE user_id = ?`

	binaryUUID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(deleteSQL, binaryUUID)
	return err
}

type sqlxUser struct {
	UserID   uuid.UUID `db:"user_id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     int       `db:"role"`
}
