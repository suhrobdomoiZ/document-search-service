package handler

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/service"
)

type documentsHandler struct {
	service *service.DocumentsService
	logger  *slog.Logger
}

func newDocumentsHandler(logger *slog.Logger, pool *pgxpool.Pool) *documentsHandler {
	service := service.NewDocumentsService(logger, pool)

	return &nodeHandler{service: service, logger: logger}
}
