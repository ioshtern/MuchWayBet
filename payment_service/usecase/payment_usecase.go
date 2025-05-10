package usecase

import (
	"errors"
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
	if p.Type != "deposit" && p.Type != "withdraw" {
		return errors.New("invalid payment type: must be 'deposit' or 'withdraw'")
	}

	if p.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = "pending"

	err := uc.repo.Create(p)
	if err != nil {
		return err
	}

	err = uc.repo.UpdateUserBalance(p.UserID, p.Amount, p.Type)
	if err != nil {
		p.Status = "failed"
		uc.repo.UpdateStatus(p.ID, "failed")
		return err
	}

	p.Status = "completed"
	return uc.repo.UpdateStatus(p.ID, "completed")
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

func (uc *PaymentUsecase) ProcessPayment(userID string, amount float64, paymentType string) (*domain.Payment, error) {
	payment := &domain.Payment{
		UserID: userID,
		Type:   paymentType,
		Amount: amount,
	}

	err := uc.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
