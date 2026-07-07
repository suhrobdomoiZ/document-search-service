package repository

import (
	"context"

	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/models"
)

type ISeedRepository interface {
	Count(cxt context.Context) (*models.DocumentsCountResponse, error)
	Create(
		ctx context.Context,
		doc *models.DocumentWithoutID,
	) (*models.CreateDocumentResponse, error)
}
