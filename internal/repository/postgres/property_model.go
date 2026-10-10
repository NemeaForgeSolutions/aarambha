package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// PropertyModel maps to the 'properties' table in PostgreSQL.
type PropertyModel struct {
	ID           string                      `db:"id"`
	Slug         string                      `db:"slug"`
	Title        string                      `db:"title"`
	Type         string                      `db:"type"`
	SubType      sql.NullString              `db:"sub_type"`
	Transaction  string                      `db:"transaction"`
	Price        int64                       `db:"price"`
	PriceDisplay sql.NullString              `db:"price_display"`
	PricePerSqFt sql.NullString              `db:"price_per_sq_ft"`
	Deposit      sql.NullString              `db:"deposit"`
	Bedrooms     int                         `db:"bedrooms"`
	Bathrooms    int                         `db:"bathrooms"`
	Balconies    int                         `db:"balconies"`
	Area         int                         `db:"area"`
	CarpetArea   int                         `db:"carpet_area"`
	AreaUnit     string                      `db:"area_unit"`
	Location     string                      `db:"location"`
	City         string                      `db:"city"`
	State        sql.NullString              `db:"state"`
	Facing       sql.NullString              `db:"facing"`
	Floor        sql.NullString              `db:"floor"`
	Parking      sql.NullString              `db:"parking"`
	Furnishing   sql.NullString              `db:"furnishing"`
	Possession   sql.NullString              `db:"possession"`
	PropertyAge  sql.NullString              `db:"property_age"`
	ReraStatus   sql.NullString              `db:"rera_status"`
	ReraNumber   sql.NullString              `db:"rera_number"`
	Featured     bool                        `db:"featured"`
	Verified     bool                        `db:"verified"`
	IsNewLaunch  bool                        `db:"is_new_launch"`
	Image        sql.NullString              `db:"image"`
	Gallery      JSONB[[]string]             `db:"gallery"`
	Description  sql.NullString              `db:"description"`
	Amenities    JSONB[[]string]             `db:"amenities"`
	Nearby       JSONB[[]domain.NearbyPlace] `db:"nearby"`
	Developer    sql.NullString              `db:"developer"`
	AgentID      sql.NullString              `db:"agent_id"`
	ProjectID    sql.NullString              `db:"project_id"`
	LocationID   sql.NullString              `db:"location_id"`
	CreatedBy    sql.NullString              `db:"created_by"`
	UpdatedBy    sql.NullString              `db:"updated_by"`
	CreatedAt    time.Time                   `db:"created_at"`
	UpdatedAt    time.Time                   `db:"updated_at"`
}

const PropertyColumns = `
	id, slug, title, type, sub_type, transaction, price, price_display, price_per_sq_ft,
	deposit, bedrooms, bathrooms, balconies, area, carpet_area, area_unit, location,
	city, state, facing, floor, parking, furnishing, possession, property_age,
	rera_status, rera_number, featured, verified, is_new_launch, image, gallery,
	description, amenities, nearby, developer, agent_id, project_id, location_id,
	created_by, updated_by
`

func (m *PropertyModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Slug, &m.Title, &m.Type, &m.SubType, &m.Transaction,
		&m.Price, &m.PriceDisplay, &m.PricePerSqFt, &m.Deposit,
		&m.Bedrooms, &m.Bathrooms, &m.Balconies, &m.Area, &m.CarpetArea,
		&m.AreaUnit, &m.Location, &m.City, &m.State, &m.Facing,
		&m.Floor, &m.Parking, &m.Furnishing, &m.Possession,
		&m.PropertyAge, &m.ReraStatus, &m.ReraNumber, &m.Featured,
		&m.Verified, &m.IsNewLaunch, &m.Image, &m.Gallery,
		&m.Description, &m.Amenities, &m.Nearby, &m.Developer, &m.AgentID,
		&m.ProjectID, &m.LocationID, &m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *PropertyModel) InsertValues() []any {
	return []any{
		m.ID, m.Slug, m.Title, m.Type, m.SubType, m.Transaction,
		m.Price, m.PriceDisplay, m.PricePerSqFt, m.Deposit,
		m.Bedrooms, m.Bathrooms, m.Balconies, m.Area, m.CarpetArea,
		m.AreaUnit, m.Location, m.City, m.State, m.Facing,
		m.Floor, m.Parking, m.Furnishing, m.Possession, m.PropertyAge,
		m.ReraStatus, m.ReraNumber, m.Featured, m.Verified, m.IsNewLaunch,
		m.Image, m.Gallery, m.Description, m.Amenities, m.Nearby,
		m.Developer, m.AgentID, m.ProjectID, m.LocationID, m.CreatedBy, m.UpdatedBy,
	}
}

func (m *PropertyModel) ToDomain() domain.Property {
	return domain.Property{
		ID:           m.ID,
		Slug:         m.Slug,
		Title:        m.Title,
		Type:         m.Type,
		SubType:      fromNullString(m.SubType),
		Transaction:  m.Transaction,
		Price:        m.Price,
		PriceDisplay: fromNullString(m.PriceDisplay),
		PricePerSqFt: fromNullString(m.PricePerSqFt),
		Deposit:      fromNullString(m.Deposit),
		Bedrooms:     m.Bedrooms,
		Bathrooms:    m.Bathrooms,
		Balconies:    m.Balconies,
		Area:         m.Area,
		CarpetArea:   m.CarpetArea,
		AreaUnit:     m.AreaUnit,
		Location:     m.Location,
		City:         m.City,
		State:        fromNullString(m.State),
		Facing:       fromNullString(m.Facing),
		Floor:        fromNullString(m.Floor),
		Parking:      fromNullString(m.Parking),
		Furnishing:   fromNullString(m.Furnishing),
		Possession:   fromNullString(m.Possession),
		PropertyAge:  fromNullString(m.PropertyAge),
		ReraStatus:   fromNullString(m.ReraStatus),
		ReraNumber:   fromNullString(m.ReraNumber),
		Featured:     m.Featured,
		Verified:     m.Verified,
		IsNewLaunch:  m.IsNewLaunch,
		Image:        fromNullString(m.Image),
		Gallery:      m.Gallery.Data,
		Description:  fromNullString(m.Description),
		Amenities:    m.Amenities.Data,
		Nearby:       m.Nearby.Data,
		Developer:    fromNullString(m.Developer),
		AgentID:      fromNullString(m.AgentID),
		ProjectID:    fromNullStringPointer(m.ProjectID),
		LocationID:   fromNullStringPointer(m.LocationID),
		CreatedBy:    fromNullStringPointer(m.CreatedBy),
		UpdatedBy:    fromNullStringPointer(m.UpdatedBy),
	}
}

func PropertyModelFromDomain(p *domain.Property) PropertyModel {
	return PropertyModel{
		ID:           p.ID,
		Slug:         p.Slug,
		Title:        p.Title,
		Type:         p.Type,
		SubType:      toNullString(p.SubType),
		Transaction:  p.Transaction,
		Price:        p.Price,
		PriceDisplay: toNullString(p.PriceDisplay),
		PricePerSqFt: toNullString(p.PricePerSqFt),
		Deposit:      toNullString(p.Deposit),
		Bedrooms:     p.Bedrooms,
		Bathrooms:    p.Bathrooms,
		Balconies:    p.Balconies,
		Area:         p.Area,
		CarpetArea:   p.CarpetArea,
		AreaUnit:     p.AreaUnit,
		Location:     p.Location,
		City:         p.City,
		State:        toNullString(p.State),
		Facing:       toNullString(p.Facing),
		Floor:        toNullString(p.Floor),
		Parking:      toNullString(p.Parking),
		Furnishing:   toNullString(p.Furnishing),
		Possession:   toNullString(p.Possession),
		PropertyAge:  toNullString(p.PropertyAge),
		ReraStatus:   toNullString(p.ReraStatus),
		ReraNumber:   toNullString(p.ReraNumber),
		Featured:     p.Featured,
		Verified:     p.Verified,
		IsNewLaunch:  p.IsNewLaunch,
		Image:        toNullString(p.Image),
		Gallery:      JSONB[[]string]{Data: p.Gallery},
		Description:  toNullString(p.Description),
		Amenities:    JSONB[[]string]{Data: p.Amenities},
		Nearby:       JSONB[[]domain.NearbyPlace]{Data: p.Nearby},
		Developer:    toNullString(p.Developer),
		AgentID:      toNullString(p.AgentID),
		ProjectID:    toNullStringPointer(p.ProjectID),
		LocationID:   toNullStringPointer(p.LocationID),
		CreatedBy:    toNullStringPointer(p.CreatedBy),
		UpdatedBy:    toNullStringPointer(p.UpdatedBy),
	}
}
