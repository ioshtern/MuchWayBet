package usecase

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"muchway/user_service/domain"
	"muchway/user_service/rabbitmq"
	"muchway/user_service/repository"
)

type UserUsecase interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
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
