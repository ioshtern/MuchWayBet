package postgres

import (
	"database/sql"
	"muchway/user_service/domain"
	"muchway/user_service/repository"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (username, password, email, balance, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, user.Username, user.Password, user.Email, user.Balance, user.Role)
	return err
}

func (r *PostgresUserRepository) GetByID(id int64) (*domain.User, error) {
	query := `SELECT id, username, password, email, balance, role FROM users WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, password, email, balance, role FROM users WHERE username = $1`
	row := r.DB.QueryRow(query, username)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, password, email, balance, role FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetAll() ([]*domain.User, error) {
	query := `SELECT id, username, password, email, balance, role FROM users`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	query := `UPDATE users SET password = $1, email = $2, balance = $3, role = $4 WHERE id = $5`
	_, err := r.DB.Exec(query, user.Password, user.Email, user.Balance, user.Role, user.ID)
	return err
}

func (r *PostgresUserRepository) Delete(username string) error {
	query := `DELETE FROM users WHERE username = $1`
	_, err := r.DB.Exec(query, username)
	return err
}
