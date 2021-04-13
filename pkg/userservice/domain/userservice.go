package domain

type UserService interface {
	AddUser(email, password string, role Role) (UserID, error)
}

func NewUserService(repository UserRepository) UserService {
	return &userService{
		repo: repository,
	}
}

type userService struct {
	repo UserRepository
}

func (service *userService) AddUser(email, password string, role Role) (UserID, error) {
	user, err := service.repo.FindByEmail(email)
	if err != nil {
		if err != ErrUserNotFound {
			return UserID{}, err
		}

		// user found
		if user.Email == email {
			return UserID{}, ErrUserWithEmailAlreadyExists
		}
	}

	user = User{
		ID:       service.repo.NewID(),
		Email:    email,
		Password: password,
		Role:     role,
	}
	err = service.repo.Store(user)

	if err != nil {
		return UserID{}, err
	}

	return user.ID, nil
}
