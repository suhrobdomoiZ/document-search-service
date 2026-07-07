package service

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/repository"
)

type DocumentsService struct {
	repository repository.IDocumentsRepository
	logger     *slog.Logger
}

func NewDocumentsService(logger *slog.Logger, pool *pgxpool.Pool) *DocumentsService {
	return &DocumentsService{logger: logger, repository: repository.NewDocuments(pool)}
}
