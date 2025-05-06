package usecase

import (
	"log"
	"muchway/user_service/domain"
	"muchway/user_service/rabbitmq"
	"muchway/user_service/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(username string) error
}

type userUsecase struct {
	repo      repository.UserRepository
	publisher *rabbitmq.Publisher
}

func NewUserUsecase(repo repository.UserRepository, publisher *rabbitmq.Publisher) UserUsecase {
	return &userUsecase{repo: repo, publisher: publisher}
}

func (u *userUsecase) Register(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := u.repo.Create(user); err != nil {
		return err
	}

	if u.publisher != nil {
		if err := u.publisher.Publish("user.created", user); err != nil {
			log.Println("❌ Failed to publish user.created:", err)
		}
	}

	return nil
}

func (u *userUsecase) Login(email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	if u.publisher != nil {
		if err := u.publisher.Publish("user.logged_in", user); err != nil {
			log.Println("❌ Failed to publish user.logged_in:", err)
		}
	}

	return user, nil
}

func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	return u.repo.GetAll()
}

func (u *userUsecase) GetUserByID(id int64) (*domain.User, error) {
	return u.repo.GetByID(id)
}

func (u *userUsecase) GetUserByUsername(username string) (*domain.User, error) {
	return u.repo.GetByUsername(username)
}

func (u *userUsecase) GetUserByEmail(email string) (*domain.User, error) {
	return u.repo.GetByEmail(email)
}

func (u *userUsecase) UpdateUser(user *domain.User) error {
	err := u.repo.Update(user)
	if err == nil && u.publisher != nil {
		_ = u.publisher.Publish("user.updated", user)
	}
	return err
}

func (u *userUsecase) DeleteUser(username string) error {
	err := u.repo.Delete(username)
	if err == nil && u.publisher != nil {
		_ = u.publisher.Publish("user.deleted", map[string]string{"username": username})
	}
	return err
}
