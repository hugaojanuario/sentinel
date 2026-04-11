package repository

import (
	"database/sql"
	"fmt"

	"github.com/hugaojanuario/sentinel/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(req models.CreateUserRequest, pass string) (*models.User, error) {
	query := `
        INSERT INTO users (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, name, email, created_at
    `
	user := &models.User{}
	err := r.db.QueryRow(query, req.Name, req.Email, pass).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar o novo usuario: %w", err)
	}

	return user, nil
}

func (r *Repository) FindUserByEMail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash, created_at
				FROM users
				WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("erro ao listar o usuario: %w", err)
	}

	return user, nil
}
