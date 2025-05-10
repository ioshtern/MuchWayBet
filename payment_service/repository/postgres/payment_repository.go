package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"muchway/payment_service/domain"
	"time"
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

func (r *PostgresPaymentRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *PostgresPaymentRepository) UpdateUserBalance(userID string, amount float64, operation string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	log.Printf("Updating balance for user %s: %s %.2f", userID, operation, amount)

	var userIDInt int64
	var currentBalance float64

	err = tx.QueryRow("SELECT id, balance FROM users WHERE id = $1", userID).Scan(&userIDInt, &currentBalance)
	if err != nil {
		log.Printf("Error querying by ID: %v", err)
		if err == sql.ErrNoRows {
			log.Printf("User ID %s not found, trying as username", userID)
			err = tx.QueryRow("SELECT id, balance FROM users WHERE username = $1", userID).Scan(&userIDInt, &currentBalance)
			if err != nil {
				log.Printf("Error querying by username: %v", err)
				if err == sql.ErrNoRows {
					return errors.New("user not found")
				}
				return err
			}
		} else {
			return err
		}
	}

	log.Printf("Found user %d with current balance: %.2f", userIDInt, currentBalance)

	var newBalance float64
	switch operation {
	case "deposit":
		newBalance = currentBalance + amount
		log.Printf("Deposit operation: %.2f + %.2f = %.2f", currentBalance, amount, newBalance)
	case "withdraw":
		if currentBalance < amount {
			log.Printf("Insufficient balance: %.2f < %.2f", currentBalance, amount)
			return errors.New("insufficient balance")
		}
		newBalance = currentBalance - amount
		log.Printf("Withdraw operation: %.2f - %.2f = %.2f", currentBalance, amount, newBalance)
	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}

	result, err := tx.Exec("UPDATE users SET balance = $1 WHERE id = $2", newBalance, userIDInt)
	if err != nil {
		log.Printf("Error updating balance: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Updated balance for user %d: %.2f, rows affected: %d", userIDInt, newBalance, rowsAffected)

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("Successfully updated balance for user %d: %s %.2f, new balance: %.2f",
		userIDInt, operation, amount, newBalance)
	return nil
}
