package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"muchway/payment_service/domain"
	"muchway/payment_service/repository"
	"time"
)

const (
	paymentKeyPrefix = "payment:"
	allPaymentsKey   = "payments:all"
	cacheTTL         = 24 * time.Hour
)

type RedisPaymentRepository struct {
	repo repository.PaymentRepository
}

func NewRedisPaymentRepository(repo repository.PaymentRepository) *RedisPaymentRepository {
	return &RedisPaymentRepository{
		repo: repo,
	}
}

func (r *RedisPaymentRepository) Create(payment *domain.Payment) error {
	log.Printf("Redis repository: Creating payment with ID %s", payment.ID)

	// First, create the payment in the underlying repository
	err := r.repo.Create(payment)
	if err != nil {
		log.Printf("Redis repository: Error creating payment in underlying repo: %v", err)
		return err
	}

	paymentJSON, err := json.Marshal(payment)
	if err != nil {
		log.Printf("Redis repository: Error marshaling payment: %v", err)
		return nil
	}

	key := fmt.Sprintf("%s%s", paymentKeyPrefix, payment.ID)
	log.Printf("Redis repository: Caching payment to key: %s", key)
	err = RedisClient.Set(Ctx, key, paymentJSON, cacheTTL).Err()
	if err != nil {
		log.Printf("Redis repository: Error caching payment: %v", err)
	} else {
		log.Printf("Redis repository: Successfully cached payment with ID %s", payment.ID)
	}

	r.invalidateAllPaymentsCache()

	return nil
}

func (r *RedisPaymentRepository) GetByID(id string) (*domain.Payment, error) {
	log.Printf("Redis repository: Getting payment with ID %s", id)

	key := fmt.Sprintf("%s%s", paymentKeyPrefix, id)
	log.Printf("Redis repository: Checking cache for key: %s", key)
	paymentJSON, err := RedisClient.Get(Ctx, key).Result()

	if err == nil {
		log.Printf("Redis repository: Cache HIT for payment ID %s", id)
		var payment domain.Payment
		err = json.Unmarshal([]byte(paymentJSON), &payment)
		if err == nil {
			log.Printf("Redis repository: Successfully retrieved payment %s from cache", id)
			return &payment, nil
		}
		log.Printf("Redis repository: Error unmarshaling payment from cache: %v", err)
	} else {
		log.Printf("Redis repository: Cache MISS for payment ID %s: %v", id, err)
	}

	log.Printf("Redis repository: Fetching payment %s from underlying repository", id)
	payment, err := r.repo.GetByID(id)
	if err != nil {
		log.Printf("Redis repository: Error fetching payment from underlying repo: %v", err)
		return nil, err
	}

	log.Printf("Redis repository: Caching payment %s for future requests", id)
	paymentData, marshalErr := json.Marshal(payment)
	if marshalErr == nil {
		err = RedisClient.Set(Ctx, key, paymentData, cacheTTL).Err()
		if err != nil {
			log.Printf("Redis repository: Error caching payment: %v", err)
		} else {
			log.Printf("Redis repository: Successfully cached payment %s", id)
		}
	} else {
		log.Printf("Redis repository: Error marshaling payment for cache: %v", marshalErr)
	}

	return payment, nil
}

func (r *RedisPaymentRepository) GetAll() ([]*domain.Payment, error) {
	paymentsJSON, err := RedisClient.Get(Ctx, allPaymentsKey).Result()

	if err == nil {
		var payments []*domain.Payment
		err = json.Unmarshal([]byte(paymentsJSON), &payments)
		if err == nil {
			return payments, nil
		}
		log.Printf("Error unmarshaling payments from cache: %v", err)
	}

	payments, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	paymentsData, marshalErr := json.Marshal(payments)
	if marshalErr == nil {
		RedisClient.Set(Ctx, allPaymentsKey, paymentsData, cacheTTL)
	}

	return payments, nil
}

func (r *RedisPaymentRepository) DeleteByID(id string) error {
	err := r.repo.DeleteByID(id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%s", paymentKeyPrefix, id)
	RedisClient.Del(Ctx, key)

	r.invalidateAllPaymentsCache()

	return nil
}

func (r *RedisPaymentRepository) UpdateStatus(id string, status string) error {
	err := r.repo.UpdateStatus(id, status)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%s", paymentKeyPrefix, id)
	RedisClient.Del(Ctx, key)

	r.invalidateAllPaymentsCache()

	return nil
}

func (r *RedisPaymentRepository) UpdateUserBalance(userID string, amount float64, operation string) error {
	err := r.repo.UpdateUserBalance(userID, amount, operation)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisPaymentRepository) invalidateAllPaymentsCache() {
	log.Printf("Invalidating: %s", allPaymentsKey)
	err := RedisClient.Del(Ctx, allPaymentsKey).Err()
	if err != nil {
		log.Printf("Error invalidating: %v", err)
	} else {
		log.Printf("Redis repository: invalidated")
	}
}
