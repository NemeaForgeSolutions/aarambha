package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// LeadModel maps to the 'leads' table in PostgreSQL.
type LeadModel struct {
	ID            string         `db:"id"`
	Name          string         `db:"name"`
	Phone         string         `db:"phone"`
	Email         sql.NullString `db:"email"`
	PropertyID    sql.NullString `db:"property_id"`
	PropertyTitle sql.NullString `db:"property_title"`
	Message       sql.NullString `db:"message"`
	Type          sql.NullString `db:"type"`
	Status        sql.NullString `db:"status"`
	PreferredDate sql.NullString `db:"preferred_date"`
	AgentID       sql.NullString `db:"agent_id"`
	UserID        sql.NullString `db:"user_id"`
	CreatedBy     sql.NullString `db:"created_by"`
	UpdatedBy     sql.NullString `db:"updated_by"`
	CreatedAt     time.Time      `db:"created_at"`
}

const LeadColumns = `
	id, name, phone, email, property_id, property_title, message, type, status, preferred_date,
	agent_id, user_id, created_by, updated_by, created_at
`

func (m *LeadModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Name, &m.Phone, &m.Email,
		&m.PropertyID, &m.PropertyTitle, &m.Message,
		&m.Type, &m.Status, &m.PreferredDate,
		&m.AgentID, &m.UserID, &m.CreatedBy, &m.UpdatedBy,
		&m.CreatedAt,
	}
}

func (m *LeadModel) InsertValues() []any {
	return []any{
		m.ID, m.Name, m.Phone, m.Email,
		m.PropertyID, m.PropertyTitle, m.Message,
		m.Type, m.Status, m.PreferredDate,
		m.AgentID, m.UserID, m.CreatedBy, m.UpdatedBy,
		m.CreatedAt,
	}
}

func (m *LeadModel) ToDomain() domain.Lead {
	st := fromNullString(m.Status)
	if st == "" {
		st = "pending"
	}
	return domain.Lead{
		ID:            m.ID,
		Name:          m.Name,
		Phone:         m.Phone,
		Email:         fromNullString(m.Email),
		PropertyID:    fromNullString(m.PropertyID),
		PropertyTitle: fromNullString(m.PropertyTitle),
		Message:       fromNullString(m.Message),
		Type:          fromNullString(m.Type),
		Status:        st,
		PreferredDate: fromNullString(m.PreferredDate),
		AgentID:       fromNullStringPointer(m.AgentID),
		UserID:        fromNullStringPointer(m.UserID),
		CreatedBy:     fromNullStringPointer(m.CreatedBy),
		UpdatedBy:     fromNullStringPointer(m.UpdatedBy),
		CreatedAt:     m.CreatedAt,
	}
}

func LeadModelFromDomain(l *domain.Lead) LeadModel {
	st := l.Status
	if st == "" {
		st = "pending"
	}
	return LeadModel{
		ID:            l.ID,
		Name:          l.Name,
		Phone:         l.Phone,
		Email:         toNullString(l.Email),
		PropertyID:    toNullString(l.PropertyID),
		PropertyTitle: toNullString(l.PropertyTitle),
		Message:       toNullString(l.Message),
		Type:          toNullString(l.Type),
		Status:        toNullString(st),
		PreferredDate: toNullString(l.PreferredDate),
		AgentID:       toNullStringPointer(l.AgentID),
		UserID:        toNullStringPointer(l.UserID),
		CreatedBy:     toNullStringPointer(l.CreatedBy),
		UpdatedBy:     toNullStringPointer(l.UpdatedBy),
		CreatedAt:     l.CreatedAt,
	}
}
