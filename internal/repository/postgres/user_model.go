package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// UserModel maps to the 'users' table in PostgreSQL.
type UserModel struct {
	ID           string         `db:"id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	PasswordHash string         `db:"password_hash"`
	Role         string         `db:"role"`
	Phone        sql.NullString `db:"phone"`
	Avatar       sql.NullString `db:"avatar"`
	IsActive     bool           `db:"is_active"`
	LastLoginAt  sql.NullTime   `db:"last_login_at"`
	CreatedBy    sql.NullString `db:"created_by"`
	UpdatedBy    sql.NullString `db:"updated_by"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

const UserColumns = `
	id, name, email, password_hash, role, phone, avatar, is_active, last_login_at, created_by, updated_by, created_at, updated_at
`

func (m *UserModel) ScanValues() []any {
	return []any{
		&m.ID,
		&m.Name,
		&m.Email,
		&m.PasswordHash,
		&m.Role,
		&m.Phone,
		&m.Avatar,
		&m.IsActive,
		&m.LastLoginAt,
		&m.CreatedBy,
		&m.UpdatedBy,
		&m.CreatedAt,
		&m.UpdatedAt,
	}
}

func (m *UserModel) ToDomain() domain.User {
	var lastLogin *time.Time
	if m.LastLoginAt.Valid {
		lastLogin = &m.LastLoginAt.Time
	}

	return domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
		Phone:        fromNullString(m.Phone),
		Avatar:       fromNullString(m.Avatar),
		IsActive:     m.IsActive,
		LastLoginAt:  lastLogin,
		CreatedBy:    fromNullStringPointer(m.CreatedBy),
		UpdatedBy:    fromNullStringPointer(m.UpdatedBy),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func UserModelFromDomain(u *domain.User) UserModel {
	var lastLogin sql.NullTime
	if u.LastLoginAt != nil {
		lastLogin = sql.NullTime{Time: *u.LastLoginAt, Valid: true}
	}

	return UserModel{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		Phone:        toNullString(u.Phone),
		Avatar:       toNullString(u.Avatar),
		IsActive:     u.IsActive,
		LastLoginAt:  lastLogin,
		CreatedBy:    toNullStringPointer(u.CreatedBy),
		UpdatedBy:    toNullStringPointer(u.UpdatedBy),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
