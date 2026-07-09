package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/document-search-service/internal/models"
)

type Seed struct {
	pool *pgxpool.Pool
}

func NewSeed(pool *pgxpool.Pool) *Seed {
	return &Seed{pool: pool}
}

func (r *Seed) Count(cxt context.Context) (*models.DocumentsCountResponse, error) {
	query := `
		SELECT COUNT(*)
		FROM documents
	`

	var result models.DocumentsCountResponse

	err := r.pool.QueryRow(cxt, query).Scan(
		&result.Count,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *Seed) Create(
	ctx context.Context,
	doc *models.DocumentWithoutID,
) (*models.CreateDocumentResponse, error) {
	query := `
		INSERT INTO documents (text, rubrics, created_date)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int

	err := r.pool.QueryRow(
		ctx,
		query,
		doc.Text,
		doc.Rubrics,
		doc.CreatedDate).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	result := &models.CreateDocumentResponse{ID: id}

	return result, nil
}
