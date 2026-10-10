package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// AgentModel maps to the 'agents' table in PostgreSQL.
type AgentModel struct {
	ID             string          `db:"id"`
	Name           string          `db:"name"`
	Role           sql.NullString  `db:"role"`
	City           sql.NullString  `db:"city"`
	AreasServed    JSONB[[]string] `db:"areas_served"`
	Specialization sql.NullString  `db:"specialization"`
	Experience     sql.NullString  `db:"experience"`
	Languages      JSONB[[]string] `db:"languages"`
	ReraID         sql.NullString  `db:"rera_id"`
	Photo          sql.NullString  `db:"photo"`
	Phone          sql.NullString  `db:"phone"`
	WhatsApp       sql.NullString  `db:"whatsapp"`
	Email          sql.NullString  `db:"email"`
	ActiveListings int             `db:"active_listings"`
	Bio            sql.NullString  `db:"bio"`
	UserID         sql.NullString  `db:"user_id"`
	CreatedBy      sql.NullString  `db:"created_by"`
	UpdatedBy      sql.NullString  `db:"updated_by"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

const AgentColumns = `
	id, name, role, city, areas_served, specialization, experience, languages,
	rera_id, photo, phone, whatsapp, email, active_listings, bio,
	user_id, created_by, updated_by
`

func (m *AgentModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Name, &m.Role, &m.City, &m.AreasServed,
		&m.Specialization, &m.Experience, &m.Languages,
		&m.ReraID, &m.Photo, &m.Phone, &m.WhatsApp, &m.Email,
		&m.ActiveListings, &m.Bio,
		&m.UserID, &m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *AgentModel) InsertValues() []any {
	return []any{
		m.ID, m.Name, m.Role, m.City, m.AreasServed,
		m.Specialization, m.Experience, m.Languages,
		m.ReraID, m.Photo, m.Phone, m.WhatsApp, m.Email,
		m.ActiveListings, m.Bio,
		m.UserID, m.CreatedBy, m.UpdatedBy,
	}
}

func (m *AgentModel) ToDomain() domain.Agent {
	return domain.Agent{
		ID:             m.ID,
		Name:           m.Name,
		Role:           fromNullString(m.Role),
		City:           fromNullString(m.City),
		AreasServed:    m.AreasServed.Data,
		Specialization: fromNullString(m.Specialization),
		Experience:     fromNullString(m.Experience),
		Languages:      m.Languages.Data,
		ReraID:         fromNullString(m.ReraID),
		Photo:          fromNullString(m.Photo),
		Phone:          fromNullString(m.Phone),
		WhatsApp:       fromNullString(m.WhatsApp),
		Email:          fromNullString(m.Email),
		ActiveListings: m.ActiveListings,
		Bio:            fromNullString(m.Bio),
		UserID:         fromNullStringPointer(m.UserID),
		CreatedBy:      fromNullStringPointer(m.CreatedBy),
		UpdatedBy:      fromNullStringPointer(m.UpdatedBy),
	}
}

func AgentModelFromDomain(a *domain.Agent) AgentModel {
	return AgentModel{
		ID:             a.ID,
		Name:           a.Name,
		Role:           toNullString(a.Role),
		City:           toNullString(a.City),
		AreasServed:    JSONB[[]string]{Data: a.AreasServed},
		Specialization: toNullString(a.Specialization),
		Experience:     toNullString(a.Experience),
		Languages:      JSONB[[]string]{Data: a.Languages},
		ReraID:         toNullString(a.ReraID),
		Photo:          toNullString(a.Photo),
		Phone:          toNullString(a.Phone),
		WhatsApp:       toNullString(a.WhatsApp),
		Email:          toNullString(a.Email),
		ActiveListings: a.ActiveListings,
		Bio:            toNullString(a.Bio),
		UserID:         toNullStringPointer(a.UserID),
		CreatedBy:      toNullStringPointer(a.CreatedBy),
		UpdatedBy:      toNullStringPointer(a.UpdatedBy),
	}
}
