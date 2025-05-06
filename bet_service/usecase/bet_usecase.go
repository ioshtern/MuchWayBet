package usecase

import (
	"bet_service/domain"
	"bet_service/repository"
	"time"
)

type BetUsecase struct {
	betRepo repository.BetRepository
}

func NewBetUsecase(betRepo repository.BetRepository) *BetUsecase {
	return &BetUsecase{betRepo: betRepo}
}

func (u *BetUsecase) CreateBet(bet *domain.Bet) error {
	now := time.Now()
	bet.CreatedAt = now
	bet.UpdatedAt = now
	bet.Status = "pending" 
	return u.betRepo.Create(bet)
}

func (u *BetUsecase) GetBetByID(id string) (*domain.Bet, error) {
	return u.betRepo.GetByID(id)
}

func (u *BetUsecase) GetBetsByUserID(userID string) ([]*domain.Bet, error) {
	return u.betRepo.GetByUserID(userID)
}

func (u *BetUsecase) UpdateBet(bet *domain.Bet) error {
	bet.UpdatedAt = time.Now()
	return u.betRepo.Update(bet)
}

func (u *BetUsecase) DeleteBet(id string) error {
	return u.betRepo.Delete(id)
}
