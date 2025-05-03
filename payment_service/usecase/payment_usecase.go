package usecase

import (
	"muchway/payment_service/domain"
	"muchway/payment_service/repository"
	"time"

	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(r repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: r}
}

func (uc *PaymentUsecase) CreatePayment(p *domain.Payment) error {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return uc.repo.Create(p)
}

func (uc *PaymentUsecase) GetPaymentByID(id string) (*domain.Payment, error) {
	return uc.repo.GetByID(id)
}

func (uc *PaymentUsecase) GetAllPayments() ([]*domain.Payment, error) {
	return uc.repo.GetAll()
}

func (uc *PaymentUsecase) DeletePaymentByID(id string) error {
	return uc.repo.DeleteByID(id)
}
