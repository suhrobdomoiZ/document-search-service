package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/document-search-service/config"
	"github.com/suhrobdomoiZ/document-search-service/internal/es"
	"github.com/suhrobdomoiZ/document-search-service/internal/repository"
	"github.com/suhrobdomoiZ/document-search-service/internal/seed"
	"github.com/suhrobdomoiZ/document-search-service/internal/server"
	"github.com/suhrobdomoiZ/document-search-service/internal/utils"
	"github.com/suhrobdomoiZ/document-search-service/migrations"
	"github.com/suhrobdomoiZ/document-search-service/pkg/closer"
	"github.com/suhrobdomoiZ/document-search-service/pkg/logger"
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

	appLogger.Info("main: ensuring es index exists...", "index", utils.EsIndexName)

	err = esClient.CreateIndex(ctx, utils.EsIndexName)
	if err != nil {
		appLogger.Error("main: failed to create es index", "error", err)
		os.Exit(1)
	}

	seedRepo := repository.NewSeed(pool)

	appLogger.Info("main: starting seeder...")

	err = seed.Run(
		ctx,
		appLogger,
		seedRepo,
		esClient,
		appConfig.SeedConfig.DataPath(),
	)
	if err != nil {
		appLogger.Error("main: seeding failed", "error", err)
		os.Exit(1)
	}

	appLogger.Info("main: seeder finished successfully")

	srv := server.NewServer(appConfig, appLogger, appCloser, pool, esClient)

	err = srv.Start(ctx)
	if err != nil {
		srv.Logger.Error("main: failed to serve", "error", err.Error())
		os.Exit(1)
	}
}
