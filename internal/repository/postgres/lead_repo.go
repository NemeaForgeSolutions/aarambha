package postgres

import (
	"context"
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
	"aarambha-agency/internal/pkg/uuid"
)

type LeadRepo struct {
	db *sql.DB
}

func NewLeadRepo(db *sql.DB) *LeadRepo {
	return &LeadRepo{db: db}
}

func (r *LeadRepo) Save(lead *domain.Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if lead.ID == "" {
		lead.ID = uuid.NewV7()
	}
	if lead.CreatedAt.IsZero() {
		lead.CreatedAt = time.Now()
	}

	m := LeadModelFromDomain(lead)

	query := `INSERT INTO leads (` + LeadColumns + `) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
	)`
	_, err := r.db.ExecContext(ctx, query, m.InsertValues()...)
	return err
}

func (r *LeadRepo) UpdateStatus(id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE leads SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *LeadRepo) FindByID(id string) (*domain.Lead, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + LeadColumns + ` FROM leads WHERE id = $1`
	var m LeadModel
	if err := r.db.QueryRowContext(ctx, query, id).Scan(m.ScanValues()...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	lead := m.ToDomain()
	return &lead, nil
}

func (r *LeadRepo) FindAll() []domain.Lead {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT ` + LeadColumns + ` FROM leads ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var list []domain.Lead
	for rows.Next() {
		var m LeadModel
		if err := rows.Scan(m.ScanValues()...); err == nil {
			list = append(list, m.ToDomain())
		}
	}
	return list
}
