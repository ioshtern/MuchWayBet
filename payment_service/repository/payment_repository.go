package repository

import "muchway/payment_service/domain"

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	GetByID(id string) (*domain.Payment, error)
	GetAll() ([]*domain.Payment, error)
	DeleteByID(id string) error
}
