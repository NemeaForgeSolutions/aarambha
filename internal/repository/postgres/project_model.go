package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// ProjectModel maps to the 'projects' table in PostgreSQL.
type ProjectModel struct {
	ID             string          `db:"id"`
	Name           string          `db:"name"`
	Developer      sql.NullString  `db:"developer"`
	Location       sql.NullString  `db:"location"`
	LocationID     sql.NullString  `db:"location_id"`
	LeadAgentID    sql.NullString  `db:"lead_agent_id"`
	StartingPrice  sql.NullString  `db:"starting_price"`
	Configuration  sql.NullString  `db:"configuration"`
	PossessionDate sql.NullString  `db:"possession_date"`
	Status         sql.NullString  `db:"status"`
	Badge          sql.NullString  `db:"badge"`
	TotalUnits     sql.NullString  `db:"total_units"`
	LandArea       sql.NullString  `db:"land_area"`
	ReraNumber     sql.NullString  `db:"rera_number"`
	Image          sql.NullString  `db:"image"`
	Overview       sql.NullString  `db:"overview"`
	Highlights     JSONB[[]string] `db:"highlights"`
	CreatedBy      sql.NullString  `db:"created_by"`
	UpdatedBy      sql.NullString  `db:"updated_by"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

const ProjectColumns = `
	id, name, developer, location, location_id, lead_agent_id, starting_price, configuration, possession_date,
	status, badge, total_units, land_area, rera_number, image, overview, highlights, created_by, updated_by
`

func (m *ProjectModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Name, &m.Developer, &m.Location, &m.LocationID, &m.LeadAgentID,
		&m.StartingPrice, &m.Configuration, &m.PossessionDate, &m.Status,
		&m.Badge, &m.TotalUnits, &m.LandArea, &m.ReraNumber,
		&m.Image, &m.Overview, &m.Highlights,
		&m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *ProjectModel) ToDomain() domain.Project {
	return domain.Project{
		ID:             m.ID,
		Name:           m.Name,
		Developer:      fromNullString(m.Developer),
		Location:       fromNullString(m.Location),
		LocationID:     fromNullStringPointer(m.LocationID),
		LeadAgentID:    fromNullStringPointer(m.LeadAgentID),
		StartingPrice:  fromNullString(m.StartingPrice),
		Configuration:  fromNullString(m.Configuration),
		PossessionDate: fromNullString(m.PossessionDate),
		Status:         fromNullString(m.Status),
		Badge:          fromNullString(m.Badge),
		TotalUnits:     fromNullString(m.TotalUnits),
		LandArea:       fromNullString(m.LandArea),
		ReraNumber:     fromNullString(m.ReraNumber),
		Image:          fromNullString(m.Image),
		Overview:       fromNullString(m.Overview),
		Highlights:     m.Highlights.Data,
		CreatedBy:      fromNullStringPointer(m.CreatedBy),
		UpdatedBy:      fromNullStringPointer(m.UpdatedBy),
	}
}

func ProjectModelFromDomain(p *domain.Project) ProjectModel {
	return ProjectModel{
		ID:             p.ID,
		Name:           p.Name,
		Developer:      toNullString(p.Developer),
		Location:       toNullString(p.Location),
		LocationID:     toNullStringPointer(p.LocationID),
		LeadAgentID:    toNullStringPointer(p.LeadAgentID),
		StartingPrice:  toNullString(p.StartingPrice),
		Configuration:  toNullString(p.Configuration),
		PossessionDate: toNullString(p.PossessionDate),
		Status:         toNullString(p.Status),
		Badge:          toNullString(p.Badge),
		TotalUnits:     toNullString(p.TotalUnits),
		LandArea:       toNullString(p.LandArea),
		ReraNumber:     toNullString(p.ReraNumber),
		Image:          toNullString(p.Image),
		Overview:       toNullString(p.Overview),
		Highlights:     JSONB[[]string]{Data: p.Highlights},
		CreatedBy:      toNullStringPointer(p.CreatedBy),
		UpdatedBy:      toNullStringPointer(p.UpdatedBy),
	}
}
