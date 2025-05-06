package usecase

import (
	"bet_service/domain"
	"bet_service/repository"
	"time"
)

type BetUsecase struct {
	betRepo   repository.BetRepository
	publisher domain.BetEventPublisher
}

func NewBetUsecase(betRepo repository.BetRepository, publisher domain.BetEventPublisher) *BetUsecase {
	return &BetUsecase{betRepo: betRepo, publisher: publisher}
}

func (u *BetUsecase) CreateBet(bet *domain.Bet) error {
	now := time.Now()
	bet.CreatedAt = now
	bet.UpdatedAt = now
	bet.Status = "pending"
	if err := u.betRepo.Create(bet); err != nil {
		return err
	}
	return u.publisher.PublishBetCreated(bet)
}

func (u *BetUsecase) UpdateBet(bet *domain.Bet) error {
	bet.UpdatedAt = time.Now()
	if err := u.betRepo.Update(bet); err != nil {
		return err
	}
	return u.publisher.PublishBetUpdated(bet)
}

func (u *BetUsecase) DeleteBet(id string) error {
	bet, err := u.betRepo.GetByID(id)
	if err != nil {
		return err
	}
	if err := u.betRepo.Delete(id); err != nil {
		return err
	}
	return u.publisher.PublishBetDeleted(bet)
}

func (u *BetUsecase) GetBetByID(id string) (*domain.Bet, error) {
	return u.betRepo.GetByID(id)
}

func (u *BetUsecase) GetBetsByUserID(userID string) ([]*domain.Bet, error) {
	return u.betRepo.GetByUserID(userID)
}
