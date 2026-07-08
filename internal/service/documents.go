package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/es"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/models"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/repository"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/utils"
)

type DocumentsService struct {
	repository repository.IDocumentsRepository
	logger     *slog.Logger
	esClient   *es.Client
}

func NewDocumentsService(
	logger *slog.Logger,
	pool *pgxpool.Pool,
	esClient *es.Client,
) *DocumentsService {
	return &DocumentsService{
		logger:     logger,
		repository: repository.NewDocuments(pool, logger),
		esClient:   esClient,
	}
}

func (s *DocumentsService) Search(ctx context.Context, text string) ([]models.FullDocument, error) {
	ids, err := s.esClient.SearchDocuments(ctx, utils.EsIndexName, text)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []models.FullDocument{}, nil
	}

	docs, err := s.repository.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (s *DocumentsService) Delete(ctx context.Context, id int) error {
	var (
		wg           sync.WaitGroup
		dbErr, esErr error
	)

	wg.Add(1)

	wg.Go(func() {
		dbErr = s.repository.DeleteDocument(ctx, id)
	})

	go func() {
		defer wg.Done()

		esErr = s.esClient.DeleteDocument(ctx, utils.EsIndexName, id)
	}()

	wg.Wait()

	if dbErr != nil {
		return dbErr
	}

	if esErr != nil {
		return esErr
	}

	return nil
}
