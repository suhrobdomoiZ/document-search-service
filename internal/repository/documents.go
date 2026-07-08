package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/models"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/utils"
)

type Documents struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewDocuments(pool *pgxpool.Pool, logger *slog.Logger) *Documents {
	return &Documents{pool: pool, logger: logger}
}

func (r *Documents) GetByIDs(ctx context.Context, ids []int) ([]models.FullDocument, error) {
	if len(ids) == 0 {
		return []models.FullDocument{}, nil
	}

	query := `
        SELECT id, text, rubrics, created_date 
        FROM documents 
        WHERE id = ANY($1) 
        ORDER BY created_date DESC 
        LIMIT $2
    `

	rows, err := r.pool.Query(ctx, query, ids, len(ids))
	if err != nil {
		r.logger.Error("repository.GetByIDs: query failed", "query", query, "err", err)

		return nil, utils.ErrInternalServerError
	}

	defer func() {
		rows.Close()
	}()

	var docs []models.FullDocument

	for rows.Next() {
		var doc models.FullDocument

		err = rows.Scan(&doc.ID, &doc.Text, &doc.Rubrics, &doc.CreatedDate)
		if err != nil {
			r.logger.Error("repository.GetByIDs: rows.Scan failed", "query", query, "err", err)

			return nil, utils.ErrInternalServerError
		}

		docs = append(docs, doc)
	}

	err = rows.Err()
	if err != nil {
		r.logger.Error("repository.GetByIDs: rows.Err failed", "query", query, "err", err)

		return nil, utils.ErrInternalServerError
	}

	return docs, nil
}

func (r *Documents) DeleteDocument(ctx context.Context, id int) error {
	query := `DELETE FROM documents WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("repository.DeleteDocument: Exec failed", "query", query, "err", err)

		return utils.ErrInternalServerError
	}

	if tag.RowsAffected() == 0 {
		return utils.ErrNotFound
	}

	return nil
}
