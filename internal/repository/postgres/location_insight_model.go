package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// LocationInsightModel maps to the 'locations' table in PostgreSQL.
type LocationInsightModel struct {
	ID            string          `db:"id"`
	Name          string          `db:"name"`
	Tagline       sql.NullString  `db:"tagline"`
	AvgPrice      sql.NullString  `db:"avg_price"`
	RentalYield   sql.NullString  `db:"rental_yield"`
	Image         sql.NullString  `db:"image"`
	KeyLocalities JSONB[[]string] `db:"key_localities"`
	Description   sql.NullString  `db:"description"`
	CreatedBy     sql.NullString  `db:"created_by"`
	UpdatedBy     sql.NullString  `db:"updated_by"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}

const LocationInsightColumns = `
	id, name, tagline, avg_price, rental_yield, image, key_localities, description, created_by, updated_by
`

func (m *LocationInsightModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Name, &m.Tagline, &m.AvgPrice, &m.RentalYield, &m.Image,
		&m.KeyLocalities, &m.Description, &m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *LocationInsightModel) ToDomain() domain.LocationInsight {
	return domain.LocationInsight{
		ID:            m.ID,
		Name:          m.Name,
		Tagline:       fromNullString(m.Tagline),
		AvgPrice:      fromNullString(m.AvgPrice),
		RentalYield:   fromNullString(m.RentalYield),
		Image:         fromNullString(m.Image),
		KeyLocalities: m.KeyLocalities.Data,
		Description:   fromNullString(m.Description),
		CreatedBy:     fromNullStringPointer(m.CreatedBy),
		UpdatedBy:     fromNullStringPointer(m.UpdatedBy),
	}
}

func LocationInsightModelFromDomain(l *domain.LocationInsight) LocationInsightModel {
	return LocationInsightModel{
		ID:            l.ID,
		Name:          l.Name,
		Tagline:       toNullString(l.Tagline),
		AvgPrice:      toNullString(l.AvgPrice),
		RentalYield:   toNullString(l.RentalYield),
		Image:         toNullString(l.Image),
		KeyLocalities: JSONB[[]string]{Data: l.KeyLocalities},
		Description:   toNullString(l.Description),
		CreatedBy:     toNullStringPointer(l.CreatedBy),
		UpdatedBy:     toNullStringPointer(l.UpdatedBy),
	}
}
