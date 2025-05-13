package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"muchway/user_service/domain"
	"muchway/user_service/rabbitmq"
	"muchway/user_service/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error) // ← updated
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(username string) error
}

type userUsecase struct {
	repo      repository.UserRepository
	publisher *rabbitmq.Publisher
	redis     *redis.Client
}

func NewUserUsecase(repo repository.UserRepository, publisher *rabbitmq.Publisher, redis *redis.Client) UserUsecase {
	return &userUsecase{repo: repo, publisher: publisher, redis: redis}
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
			log.Println("Failed to publish user.created:", err)
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
			log.Println(" Failed to publish user.logged_in:", err)
		}
	}

	return user, nil
}

func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	return u.repo.GetAll()
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	ctx = context.Background()
	cacheKey := fmt.Sprintf("user:%d", id)

	cachedUser, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user domain.User
		if jsonErr := json.Unmarshal([]byte(cachedUser), &user); jsonErr == nil {
			fmt.Println("User fetched from cache:", user.ID)
			return &user, nil
		}
		fmt.Println("Failed to unmarshal cached user:")
	} else {
		fmt.Println("User not found in cache:", err)
	}

	user, err := u.repo.GetByID(id)
	if err != nil {
		fmt.Println("Error fetching user from database:", err)
		return nil, err
	}

	fmt.Println("User fetched from database:", user.ID)

	userData, _ := json.Marshal(user)
	if setErr := u.redis.Set(ctx, cacheKey, userData, 10*time.Minute).Err(); setErr == nil {
		fmt.Println("User cached successfully:", user.ID)
	} else {
		fmt.Println("Failed to cache user:", setErr)
	}

	return user, nil
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
