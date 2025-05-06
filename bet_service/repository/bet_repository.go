package repository

import "bet_service/domain"

type BetRepository interface {
	Create(bet *domain.Bet) error
	GetByID(id string) (*domain.Bet, error)
	GetByUserID(userID string) ([]*domain.Bet, error)
	Update(bet *domain.Bet) error
	Delete(id string) error
}
