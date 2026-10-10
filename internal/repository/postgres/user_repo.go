package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/pkg/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) scanUser(s scanner) (*domain.User, error) {
	var m UserModel
	if err := s.Scan(m.ScanValues()...); err != nil {
		return nil, err
	}
	u := m.ToDomain()
	return &u, nil
}

func (r *UserRepo) Create(u *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if u.ID == "" {
		u.ID = uuid.NewV7()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	m := UserModelFromDomain(u)

	query := `
		INSERT INTO users (id, name, email, password_hash, role, phone, avatar, is_active, last_login_at, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.Name, m.Email, m.PasswordHash, m.Role,
		m.Phone, m.Avatar, m.IsActive, m.LastLoginAt, m.CreatedBy, m.UpdatedBy, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return domain.ErrEmailAlreadyExists
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cleanedEmail := strings.ToLower(strings.TrimSpace(email))
	query := `SELECT ` + UserColumns + ` FROM users WHERE LOWER(email) = $1 LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, cleanedEmail)
	u, err := r.scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	return u, nil
}

func (r *UserRepo) FindByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + UserColumns + ` FROM users WHERE id = $1 LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, id)
	u, err := r.scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}

	return u, nil
}

func (r *UserRepo) Update(u *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u.UpdatedAt = time.Now()
	m := UserModelFromDomain(u)

	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, role = $4, phone = $5, avatar = $6, is_active = $7, updated_by = $8, updated_at = $9
		WHERE id = $10
	`

	res, err := r.db.ExecContext(ctx, query,
		m.Name, m.Email, m.PasswordHash, m.Role, m.Phone, m.Avatar, m.IsActive, m.UpdatedBy, m.UpdatedAt, m.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) UpdateLastLogin(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, now, id)
	return err
}

func (r *UserRepo) Count() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (r *UserRepo) FindAll() []domain.User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + UserColumns + ` FROM users ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err == nil {
			users = append(users, *u)
		}
	}
	return users
}
