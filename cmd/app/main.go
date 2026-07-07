package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/config"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/es"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/migrations"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/pkg/closer"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	appConfig := config.NewAppConfig()
	appLogger := logger.With("env", appConfig.EnvType())
	logger.Setup(appConfig.EnvType())

	appCloser := closer.New(appLogger)
	appCloser.AddFunc("context", cancel)

	pool, err := pgxpool.New(ctx, appConfig.DbConfig.DSN())
	if err != nil {
		appLogger.Error("main: failed to create pool", "error", err)
		os.Exit(1)
	}

	appCloser.AddFunc("pool", pool.Close)

	for attempt := range 10 {
		err = pool.Ping(ctx)
		if err == nil {
			break
		}

		appLogger.Warn("main: db is not ready", "attempt", attempt+1, "error", err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		appLogger.Error("main: db ping failed", "error", err)
		os.Exit(1)
	}

	appLogger.Info("main: connected to db", "host", appConfig.DbConfig.DBHost())

	err = migrations.Up(ctx, pool, appLogger)
	if err != nil {
		appLogger.Error("main: migration failed", "error", err)
		os.Exit(1)
	}

	appLogger.Info("main: connecting to elasticsearch...")

	esClient, err := es.NewClient(appConfig.EsConfig.Address(), appConfig.EsConfig.SearchSize())
	if err != nil {
		appLogger.Error("main: failed to create es client", "error", err)
		os.Exit(1)
	}

	appCloser.Add("es", esClient.Close)

	appLogger.Info("main: connected to elasticsearch", "address", appConfig.EsConfig.Address())
}
