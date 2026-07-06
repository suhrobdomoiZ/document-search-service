package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/*up.sql
var MigrateFS embed.FS

func Up(ctx context.Context, pool *pgxpool.Pool, appLogger *slog.Logger) error {
	entries, err := fs.ReadDir(MigrateFS, "sql")
	if err != nil {
		return fmt.Errorf("migrations.Up: fs.ReadDir failed: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		path := "sql/" + entry.Name()

		sqlBytes, err := MigrateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("migrations.Up: fs.ReadFile(%s): %w", path, err)
		}

		_, err = pool.Exec(ctx, string(sqlBytes))
		if err != nil {
			return fmt.Errorf("migrations.Up: failed to execute(pool.Exec(%s): %w)", path, err)
		}

		appLogger.Info("migrate.Up: Migration is made", "file", entry.Name())
	}

	return nil
}
