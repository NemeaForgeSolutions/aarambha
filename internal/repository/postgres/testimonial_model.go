package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// TestimonialModel maps to the 'testimonials' table in PostgreSQL.
type TestimonialModel struct {
	ID         string         `db:"id"`
	Quote      string         `db:"quote"`
	Client     string         `db:"client"`
	Location   sql.NullString `db:"location"`
	Type       sql.NullString `db:"type"`
	Property   sql.NullString `db:"property"`
	PropertyID sql.NullString `db:"property_id"`
	AgentID    sql.NullString `db:"agent_id"`
	CreatedBy  sql.NullString `db:"created_by"`
	UpdatedBy  sql.NullString `db:"updated_by"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}

const TestimonialColumns = `
	id, quote, client, location, type, property, property_id, agent_id, created_by, updated_by
`

func (m *TestimonialModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Quote, &m.Client, &m.Location, &m.Type, &m.Property,
		&m.PropertyID, &m.AgentID, &m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *TestimonialModel) ToDomain() domain.Testimonial {
	return domain.Testimonial{
		ID:         m.ID,
		Quote:      m.Quote,
		Client:     m.Client,
		Location:   fromNullString(m.Location),
		Type:       fromNullString(m.Type),
		Property:   fromNullString(m.Property),
		PropertyID: fromNullStringPointer(m.PropertyID),
		AgentID:    fromNullStringPointer(m.AgentID),
		CreatedBy:  fromNullStringPointer(m.CreatedBy),
		UpdatedBy:  fromNullStringPointer(m.UpdatedBy),
	}
}

func TestimonialModelFromDomain(t *domain.Testimonial) TestimonialModel {
	return TestimonialModel{
		ID:         t.ID,
		Quote:      t.Quote,
		Client:     t.Client,
		Location:   toNullString(t.Location),
		Type:       toNullString(t.Type),
		Property:   toNullString(t.Property),
		PropertyID: toNullStringPointer(t.PropertyID),
		AgentID:    toNullStringPointer(t.AgentID),
		CreatedBy:  toNullStringPointer(t.CreatedBy),
		UpdatedBy:  toNullStringPointer(t.UpdatedBy),
	}
}
