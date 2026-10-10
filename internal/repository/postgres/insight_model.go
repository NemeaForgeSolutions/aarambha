package postgres

import (
	"database/sql"
	"time"

	"aarambha-agency/internal/domain"
)

// InsightArticleModel maps to the 'insights' table in PostgreSQL.
type InsightArticleModel struct {
	ID            string         `db:"id"`
	Slug          string         `db:"slug"`
	Title         string         `db:"title"`
	Category      sql.NullString `db:"category"`
	ReadTime      sql.NullString `db:"read_time"`
	Date          sql.NullString `db:"date"`
	Summary       sql.NullString `db:"summary"`
	Author        sql.NullString `db:"author"`
	AuthorAgentID sql.NullString `db:"author_agent_id"`
	KeyTakeaway   sql.NullString `db:"key_takeaway"`
	CreatedBy     sql.NullString `db:"created_by"`
	UpdatedBy     sql.NullString `db:"updated_by"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

const InsightArticleColumns = `
	id, slug, title, category, read_time, date, summary, author, author_agent_id, key_takeaway, created_by, updated_by
`

func (m *InsightArticleModel) ScanValues() []any {
	return []any{
		&m.ID, &m.Slug, &m.Title, &m.Category, &m.ReadTime, &m.Date,
		&m.Summary, &m.Author, &m.AuthorAgentID, &m.KeyTakeaway,
		&m.CreatedBy, &m.UpdatedBy,
	}
}

func (m *InsightArticleModel) ToDomain() domain.InsightArticle {
	return domain.InsightArticle{
		ID:            m.ID,
		Slug:          m.Slug,
		Title:         m.Title,
		Category:      fromNullString(m.Category),
		ReadTime:      fromNullString(m.ReadTime),
		Date:          fromNullString(m.Date),
		Summary:       fromNullString(m.Summary),
		Author:        fromNullString(m.Author),
		AuthorAgentID: fromNullStringPointer(m.AuthorAgentID),
		KeyTakeaway:   fromNullString(m.KeyTakeaway),
		CreatedBy:     fromNullStringPointer(m.CreatedBy),
		UpdatedBy:     fromNullStringPointer(m.UpdatedBy),
	}
}

func InsightArticleModelFromDomain(i *domain.InsightArticle) InsightArticleModel {
	return InsightArticleModel{
		ID:            i.ID,
		Slug:          i.Slug,
		Title:         i.Title,
		Category:      toNullString(i.Category),
		ReadTime:      toNullString(i.ReadTime),
		Date:          toNullString(i.Date),
		Summary:       toNullString(i.Summary),
		Author:        toNullString(i.Author),
		AuthorAgentID: toNullStringPointer(i.AuthorAgentID),
		KeyTakeaway:   toNullString(i.KeyTakeaway),
		CreatedBy:     toNullStringPointer(i.CreatedBy),
		UpdatedBy:     toNullStringPointer(i.UpdatedBy),
	}
}
