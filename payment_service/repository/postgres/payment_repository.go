package postgres

import (
	"database/sql"
	"muchway/payment_service/domain"
)

type PostgresPaymentRepository struct {
	db *sql.DB
}

func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

func (r *PostgresPaymentRepository) Create(p *domain.Payment) error {
	query := `INSERT INTO payments (id, user_id, type, amount, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, p.ID, p.UserID, p.Type, p.Amount, p.Status, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *PostgresPaymentRepository) GetByID(id string) (*domain.Payment, error) {
	query := `SELECT id, user_id, type, amount, status, created_at, updated_at FROM payments WHERE id = $1`
	row := r.db.QueryRow(query, id)
	p := &domain.Payment{}
	err := row.Scan(&p.ID, &p.UserID, &p.Type, &p.Amount, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresPaymentRepository) GetAll() ([]*domain.Payment, error) {
	query := `SELECT id, user_id, type, amount, status, created_at, updated_at FROM payments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		p := &domain.Payment{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Type, &p.Amount, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PostgresPaymentRepository) DeleteByID(id string) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
