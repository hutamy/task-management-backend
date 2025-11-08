package repository

import (
	"database/sql"
	"fmt"
	"task-management-backend/internal/domain/entity"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE username = ?
	`

	var user entity.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(user *entity.User) error {
	query := `
		INSERT INTO users (username, password, created_at)
		VALUES (?, ?, ?)
	`

	now := time.Now()
	user.CreatedAt = now

	result, err := r.db.Exec(query, user.Username, user.Password, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return nil
}

func (r *UserRepository) Upsert(user *entity.User) (*entity.User, error) {
	existingUser, err := r.GetByUsername(user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return existingUser, nil
	}

	if err := r.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
