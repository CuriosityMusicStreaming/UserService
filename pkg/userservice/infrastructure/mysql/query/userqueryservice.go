package query

import (
	"database/sql"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"userservice/pkg/userservice/app/query"
	"userservice/pkg/userservice/domain"
)

func NewUserQueryService(client mysql.Client) query.UserQueryService {
	return &userQueryService{
		client: client,
	}
}

type userQueryService struct {
	client mysql.Client
}

func (service *userQueryService) GetUser(id uuid.UUID) (query.UserView, error) {
	const selectSQL = `SELECT * from user WHERE user_id = ?`

	binaryUUID, err := id.MarshalBinary()
	if err != nil {
		return query.UserView{}, err
	}

	var user sqlxUserView

	err = service.client.Get(&user, selectSQL, binaryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return query.UserView{}, query.ErrUserNotFound
		}
		return query.UserView{}, errors.WithStack(err)
	}

	return query.UserView{
		ID:       user.UserID,
		Email:    user.Email,
		Password: user.Password,
		Role:     query.Role(user.Role),
	}, nil
}

func (service *userQueryService) GetByEmail(email string) (query.UserView, error) {
	const selectSQL = `SELECT * from user WHERE email = ?`

	var user sqlxUserView

	err := service.client.Get(&user, selectSQL, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return query.UserView{}, domain.ErrUserNotFound
		}
		return query.UserView{}, errors.WithStack(err)
	}

	return query.UserView{
		ID:       user.UserID,
		Email:    user.Email,
		Password: user.Password,
		Role:     query.Role(user.Role),
	}, nil
}

type sqlxUserView struct {
	UserID   uuid.UUID `db:"user_id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     int       `db:"role"`
}
