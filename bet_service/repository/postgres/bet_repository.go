package postgres

import (
	"bet_service/domain"
	"database/sql"
)

type PostgresBetRepository struct {
	db *sql.DB
}

func NewPostgresBetRepository(db *sql.DB) *PostgresBetRepository {
	return &PostgresBetRepository{db: db}
}

func (r *PostgresBetRepository) Create(bet *domain.Bet) error {
	query := `
        INSERT INTO bets (id, user_id, event_id, amount, odds, status, payout, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.Exec(query,
		bet.ID,
		bet.UserID,
		bet.EventID,
		bet.Amount,
		bet.Odds,
		bet.Status,
		bet.Payout,
		bet.CreatedAt,
		bet.UpdatedAt,
	)
	return err
}

func (r *PostgresBetRepository) GetByID(id string) (*domain.Bet, error) {
	query := `
        SELECT id, user_id, event_id, amount, odds, status, payout, created_at, updated_at
        FROM bets
        WHERE id = $1
    `
	row := r.db.QueryRow(query, id)

	bet := &domain.Bet{}
	err := row.Scan(
		&bet.ID,
		&bet.UserID,
		&bet.EventID,
		&bet.Amount,
		&bet.Odds,
		&bet.Status,
		&bet.Payout,
		&bet.CreatedAt,
		&bet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return bet, nil
}

func (r *PostgresBetRepository) GetByUserID(userID string) ([]*domain.Bet, error) {
	query := `
        SELECT id, user_id, event_id, amount, odds, status, payout, created_at, updated_at
        FROM bets
        WHERE user_id = $1
    `
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []*domain.Bet
	for rows.Next() {
		bet := &domain.Bet{}
		err := rows.Scan(
			&bet.ID,
			&bet.UserID,
			&bet.EventID,
			&bet.Amount,
			&bet.Odds,
			&bet.Status,
			&bet.Payout,
			&bet.CreatedAt,
			&bet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)
	}

	return bets, nil
}

func (r *PostgresBetRepository) Update(bet *domain.Bet) error {
	query := `
        UPDATE bets
        SET user_id=$1, event_id=$2, amount=$3, odds=$4, status=$5, payout=$6, created_at=$7, updated_at=$8
        WHERE id=$9
    `
	_, err := r.db.Exec(query,
		bet.UserID,
		bet.EventID,
		bet.Amount,
		bet.Odds,
		bet.Status,
		bet.Payout,
		bet.CreatedAt,
		bet.UpdatedAt,
		bet.ID,
	)
	return err
}

func (r *PostgresBetRepository) Delete(id string) error {
	query := `
        DELETE FROM bets
        WHERE id = $1
    `
	_, err := r.db.Exec(query, id)
	return err
}
