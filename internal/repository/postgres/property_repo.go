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

type scanner interface {
	Scan(dest ...any) error
}

type PropertyRepo struct {
	db *sql.DB
}

func NewPropertyRepo(db *sql.DB) *PropertyRepo {
	return &PropertyRepo{db: db}
}

func (r *PropertyRepo) scanProperty(s scanner) (*domain.Property, error) {
	var m PropertyModel
	if err := s.Scan(m.ScanValues()...); err != nil {
		return nil, err
	}
	p := m.ToDomain()
	return &p, nil
}

func (r *PropertyRepo) FindAll() []domain.Property {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + PropertyColumns + ` FROM properties ORDER BY featured DESC, id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Property
	for rows.Next() {
		p, err := r.scanProperty(rows)
		if err == nil {
			list = append(list, *p)
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindByFilter(f domain.PropertyFilter) []domain.Property {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var conditions []string
	var args []any
	argIdx := 1

	if f.Transaction != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(transaction) = LOWER($%d)", argIdx))
		args = append(args, f.Transaction)
		argIdx++
	}
	if f.City != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(city) = LOWER($%d)", argIdx))
		args = append(args, f.City)
		argIdx++
	}
	if f.Type != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(type) = LOWER($%d)", argIdx))
		args = append(args, f.Type)
		argIdx++
	}
	if f.Bedrooms > 0 {
		conditions = append(conditions, fmt.Sprintf("bedrooms = $%d", argIdx))
		args = append(args, f.Bedrooms)
		argIdx++
	}
	if f.MinPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIdx))
		args = append(args, f.MinPrice)
		argIdx++
	}
	if f.MaxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIdx))
		args = append(args, f.MaxPrice)
		argIdx++
	}
	if f.Possession != "" {
		posPattern := "%" + strings.ToLower(f.Possession) + "%"
		conditions = append(conditions, fmt.Sprintf("LOWER(possession) LIKE $%d", argIdx))
		args = append(args, posPattern)
		argIdx++
	}
	if f.Furnishing != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(furnishing) = LOWER($%d)", argIdx))
		args = append(args, f.Furnishing)
		argIdx++
	}
	if f.ReraOnly {
		conditions = append(conditions, "LOWER(rera_status) LIKE '%rera%'")
	}
	if f.Query != "" {
		qPattern := "%" + strings.ToLower(f.Query) + "%"
		conditions = append(conditions, fmt.Sprintf(`(
			LOWER(title) LIKE $%d OR 
			LOWER(location) LIKE $%d OR 
			LOWER(city) LIKE $%d OR 
			LOWER(developer) LIKE $%d OR 
			LOWER(description) LIKE $%d
		)`, argIdx, argIdx, argIdx, argIdx, argIdx))
		args = append(args, qPattern)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	orderBy := "ORDER BY featured DESC, id DESC"
	switch f.Sort {
	case "price_asc":
		orderBy = "ORDER BY price ASC"
	case "price_desc":
		orderBy = "ORDER BY price DESC"
	case "area_desc":
		orderBy = "ORDER BY area DESC"
	}

	query := fmt.Sprintf("SELECT %s FROM properties %s %s", PropertyColumns, whereClause, orderBy)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Property
	for rows.Next() {
		p, err := r.scanProperty(rows)
		if err == nil {
			list = append(list, *p)
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindBySlug(slugOrID string) (*domain.Property, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT %s FROM properties WHERE slug = $1 OR id = $1 LIMIT 1", PropertyColumns)
	row := r.db.QueryRowContext(ctx, query, slugOrID)
	p, err := r.scanProperty(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PropertyRepo) FindFeatured(propType string, limit int) []domain.Property {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var conditions []string
	var args []any
	argIdx := 1

	conditions = append(conditions, "featured = TRUE")
	if propType != "" && !strings.EqualFold(propType, "all") {
		conditions = append(conditions, fmt.Sprintf("LOWER(type) = LOWER($%d)", argIdx))
		args = append(args, propType)
		argIdx++
	}

	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf("LIMIT %d", limit)
	}

	query := fmt.Sprintf("SELECT %s FROM properties WHERE %s ORDER BY id DESC %s",
		PropertyColumns, strings.Join(conditions, " AND "), limitClause)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Property
	for rows.Next() {
		p, err := r.scanProperty(rows)
		if err == nil {
			list = append(list, *p)
		}
	}
	if len(list) > 0 {
		return list
	}

	// Fallback: If no explicitly featured listings found for this filter, return latest verified properties
	var fallbackConditions []string
	var fallbackArgs []any
	fbArgIdx := 1
	if propType != "" && !strings.EqualFold(propType, "all") {
		fallbackConditions = append(fallbackConditions, fmt.Sprintf("LOWER(type) = LOWER($%d)", fbArgIdx))
		fallbackArgs = append(fallbackArgs, propType)
		fbArgIdx++
	}
	fbWhereClause := ""
	if len(fallbackConditions) > 0 {
		fbWhereClause = "WHERE " + strings.Join(fallbackConditions, " AND ")
	}
	fbQuery := fmt.Sprintf("SELECT %s FROM properties %s ORDER BY price DESC %s",
		PropertyColumns, fbWhereClause, limitClause)
	fbRows, fbErr := r.db.QueryContext(ctx, fbQuery, fallbackArgs...)
	if fbErr != nil {
		return nil
	}
	defer fbRows.Close()

	for fbRows.Next() {
		p, err := r.scanProperty(fbRows)
		if err == nil {
			list = append(list, *p)
		}
	}
	return list
}

func (r *PropertyRepo) FindSimilar(current *domain.Property, limit int) []domain.Property {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 3
	}

	query := fmt.Sprintf(`SELECT %s FROM properties 
		WHERE id != $1 AND (city = $2 OR type = $3) 
		ORDER BY (city = $2) DESC, price ASC LIMIT %d`, PropertyColumns, limit)

	rows, err := r.db.QueryContext(ctx, query, current.ID, current.City, current.Type)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Property
	for rows.Next() {
		p, err := r.scanProperty(rows)
		if err == nil {
			list = append(list, *p)
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) Count() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM properties").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (r *PropertyRepo) Create(prop *domain.Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if prop.ID == "" {
		prop.ID = uuid.NewV7()
	}

	m := PropertyModelFromDomain(prop)

	query := `INSERT INTO properties (` + PropertyColumns + `) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14, $15, $16, $17,
		$18, $19, $20, $21, $22, $23, $24, $25,
		$26, $27, $28, $29, $30, $31, $32,
		$33, $34, $35, $36, $37, $38, $39,
		$40, $41
	)`

	_, err := r.db.ExecContext(ctx, query, m.InsertValues()...)
	return err
}

func (r *PropertyRepo) Update(prop *domain.Property) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m := PropertyModelFromDomain(prop)

	query := `UPDATE properties SET
		slug = $2, title = $3, type = $4, sub_type = $5, transaction = $6,
		price = $7, price_display = $8, price_per_sq_ft = $9, deposit = $10,
		bedrooms = $11, bathrooms = $12, balconies = $13, area = $14,
		carpet_area = $15, area_unit = $16, location = $17, city = $18,
		state = $19, facing = $20, floor = $21, parking = $22,
		furnishing = $23, possession = $24, property_age = $25,
		rera_status = $26, rera_number = $27, featured = $28, verified = $29,
		is_new_launch = $30, image = $31, gallery = $32, description = $33,
		amenities = $34, nearby = $35, developer = $36, agent_id = $37,
		project_id = $38, location_id = $39, created_by = $40, updated_by = $41,
		updated_at = NOW()
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, m.InsertValues()...)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PropertyRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, "DELETE FROM properties WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PropertyRepo) ToggleFeatured(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newFeatured bool
	err := r.db.QueryRowContext(ctx,
		"UPDATE properties SET featured = NOT featured, updated_at = NOW() WHERE id = $1 RETURNING featured", id).
		Scan(&newFeatured)
	if errors.Is(err, sql.ErrNoRows) {
		return false, domain.ErrNotFound
	}
	if err != nil {
		return false, err
	}
	return newFeatured, nil
}

// ============================================================================
// CatalogRepository Implementation
// ============================================================================

func (r *PropertyRepo) FindAgents() []domain.Agent {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + AgentColumns + ` FROM agents ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Agent
	for rows.Next() {
		var m AgentModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindAgentByID(id string) (*domain.Agent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + AgentColumns + ` FROM agents WHERE id = $1 LIMIT 1`
	var m AgentModel
	err := r.db.QueryRowContext(ctx, query, id).Scan(m.ScanValues()...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	agent := m.ToDomain()
	return &agent, nil
}

func (r *PropertyRepo) CreateAgent(agent *domain.Agent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if agent.ID == "" {
		agent.ID = uuid.NewV7()
	}

	m := AgentModelFromDomain(agent)

	query := `INSERT INTO agents (
		id, name, role, city, areas_served, specialization, experience, languages,
		rera_id, photo, phone, whatsapp, email, active_listings, bio,
		user_id, created_by, updated_by
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8,
		$9, $10, $11, $12, $13, $14, $15,
		$16, $17, $18
	)`

	_, err := r.db.ExecContext(ctx, query, m.InsertValues()...)
	return err
}

func (r *PropertyRepo) UpdateAgent(agent *domain.Agent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m := AgentModelFromDomain(agent)

	query := `UPDATE agents SET
		name = $2, role = $3, city = $4, areas_served = $5, specialization = $6,
		experience = $7, languages = $8, rera_id = $9, photo = $10, phone = $11,
		whatsapp = $12, email = $13, active_listings = $14, bio = $15,
		user_id = $16, created_by = $17, updated_by = $18, updated_at = NOW()
		WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, m.InsertValues()...)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PropertyRepo) DeleteAgent(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(ctx, "DELETE FROM agents WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PropertyRepo) FindProjects() []domain.Project {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + ProjectColumns + ` FROM projects ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Project
	for rows.Next() {
		var m ProjectModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindLocations() []domain.LocationInsight {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + LocationInsightColumns + ` FROM locations ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.LocationInsight
	for rows.Next() {
		var m LocationInsightModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindInsights() []domain.InsightArticle {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + InsightArticleColumns + ` FROM insights ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.InsightArticle
	for rows.Next() {
		var m InsightArticleModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindTestimonials() []domain.Testimonial {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + TestimonialColumns + ` FROM testimonials ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Testimonial
	for rows.Next() {
		var m TestimonialModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}

func (r *PropertyRepo) FindWhyChooseUs() []domain.WhyChooseUsItem {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + WhyChooseUsColumns + ` FROM highlights ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.WhyChooseUsItem
	for rows.Next() {
		var m WhyChooseUsModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return list
}
