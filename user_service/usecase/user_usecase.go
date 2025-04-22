package usecase

import (
	"golang.org/x/crypto/bcrypt"
	"muchway/user_service/domain"
	"muchway/user_service/repository"
)

type UserUsecase interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Register(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return u.repo.Create(user)
}

func (u *userUsecase) Login(email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	return u.repo.GetAll()
}

func (u *userUsecase) GetUserByID(id int64) (*domain.User, error) {
	return u.repo.GetByID(id)
}
