package usecase

import (
	"errors"
	"log"
	"muchway/payment_service/client"
	"muchway/payment_service/domain"
	"muchway/payment_service/email"
	"muchway/payment_service/repository"
	"time"

	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo         repository.PaymentRepository
	userClient   *client.UserClient
	emailService email.EmailService
}

func NewPaymentUsecase(r repository.PaymentRepository, userClient *client.UserClient, emailService email.EmailService) *PaymentUsecase {
	return &PaymentUsecase{
		repo:         r,
		userClient:   userClient,
		emailService: emailService,
	}
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
	err = uc.repo.UpdateStatus(p.ID, "completed")
	if err != nil {
		return err
	}

	// Send email notification
	if uc.userClient != nil && uc.emailService != nil {
		go func() {
			// Get user email
			userEmail, err := uc.userClient.GetUserEmail(p.UserID)
			if err != nil {
				log.Printf("Failed to get user email: %v", err)
				return
			}

			// Send payment confirmation email
			err = uc.emailService.SendPaymentConfirmation(userEmail, p.UserID, p.Amount, p.Type, p.Status)
			if err != nil {
				log.Printf("Failed to send payment confirmation email: %v", err)
			} else {
				log.Printf("Payment confirmation email sent to: %s", userEmail)
			}
		}()
	}

	return nil
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
