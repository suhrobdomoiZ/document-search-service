package handler

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/es"
)

type AppHandler struct {
	DocumentsHandler *documentsHandler
}

func NewAppHandler(logger *slog.Logger, pool *pgxpool.Pool, esClient *es.Client) *AppHandler {
	return &AppHandler{
		DocumentsHandler: newDocumentsHandler(logger, pool, esClient),
	}
}
