package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// WhyChooseUsModel maps to the 'highlights' table in PostgreSQL.
type WhyChooseUsModel struct {
	ID          string         `db:"id"`
	Number      sql.NullString `db:"number"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

const WhyChooseUsColumns = `
	id, number, title, description, created_by, updated_by
`

func (m *WhyChooseUsModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Number, &m.Title, &m.Description, &m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *WhyChooseUsModel) ToDomain() domain.WhyChooseUsItem {
	return domain.WhyChooseUsItem{
		ID:          m.ID,
		Number:      fromNullString(m.Number),
		Title:       m.Title,
		Desc:        fromNullString(m.Description),
		Description: fromNullString(m.Description),
		CreatedBy:   fromNullStringPointer(m.CreatedBy),
		UpdatedBy:   fromNullStringPointer(m.UpdatedBy),
	}
}

func WhyChooseUsModelFromDomain(w *domain.WhyChooseUsItem) WhyChooseUsModel {
	desc := w.Description
	if desc == "" {
		desc = w.Desc
	}
	return WhyChooseUsModel{
		ID:          w.ID,
		Number:      toNullString(w.Number),
		Title:       w.Title,
		Description: toNullString(desc),
		CreatedBy:   toNullStringPointer(w.CreatedBy),
		UpdatedBy:   toNullStringPointer(w.UpdatedBy),
	}
}
