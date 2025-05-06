package usecase

import (
	"fmt"
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
		err := u.publisher.Publish(map[string]interface{}{
			"event": "UserCreated",
			"data": map[string]interface{}{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"balance":  user.Balance,
				"role":     user.Role,
			},
		})
		if err != nil {
			log.Println("Failed to publish UserCreated event:", err)
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

	event := fmt.Sprintf("User logged in: ID=%d, Email=%s", user.ID, user.Email)
	if err := u.publisher.Publish([]byte(event)); err != nil {
		log.Printf("Failed to publish login event: %v", err)
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
	return u.repo.Update(user)
}

func (u *userUsecase) DeleteUser(username string) error {
	return u.repo.Delete(username)
}
