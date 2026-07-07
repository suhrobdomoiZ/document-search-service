package seed

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/es"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/models"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/repository"
)

const esIndexName = "documents"

func Run(
	ctx context.Context,
	logger *slog.Logger,
	repo repository.ISeedRepository,
	esClient *es.Client,
	filePath string,
) error {
	resp, err := repo.Count(ctx)
	if err != nil {
		return fmt.Errorf("seed.Run: failed to get doc count: %w", err)
	}

	if resp.Count > 0 {
		logger.Info("seed.Run: database is not empty, skipping seeding")

		return nil
	}

	logger.Info("seed.Run: database is empty, starting to seed from CSV", "file", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("seed.Run: failed to open csv file: %w", err)
	}

	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("seed.Run: failed to read csv header: %w", err)
	}

	processedCount := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		text := record[0]

		createdAt, err := time.Parse("2006-01-02 15:04:05", record[1])
		if err != nil {
			logger.Warn("seed: failed to parse date, skipping row", "error", err)

			continue
		}

		createdAt = createdAt.UTC()

		rubrics, err := parseRubrics(record[2])
		if err != nil {
			logger.Warn("seed: failed to parse rubrics, using empty array", "error", err)

			rubrics = []string{}
		}

		doc := &models.DocumentWithoutID{
			Text:        text,
			Rubrics:     rubrics,
			CreatedDate: createdAt,
		}

		resp, err := repo.Create(ctx, doc)
		if err != nil {
			logger.Error("seed: failed to insert into db", "error", err)

			continue
		}

		err = esClient.IndexDocument(ctx, esIndexName, resp.ID, text)
		if err != nil {
			logger.Error("seed: failed to index in es", "id", resp.ID, "error", err)

			continue
		}

		processedCount++
		if processedCount%100 == 0 {
			logger.Info("seed: processing rows...", "count", processedCount)
		}
	}

	logger.Info("seed: finished successfully", "total_rows", processedCount)

	return nil
}

func parseRubrics(raw string) ([]string, error) {
	fixedJSON := strings.ReplaceAll(raw, "'", "\"")

	var rubrics []string

	err := json.Unmarshal([]byte(fixedJSON), &rubrics)
	if err != nil {
		return nil, err
	}

	return rubrics, nil
}
