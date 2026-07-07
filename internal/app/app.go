package app

import (
	"log/slog"

	"github.com/suhrobdomoiZ/anal-prog-decisions-test/config"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/internal/handler"
	"github.com/suhrobdomoiZ/anal-prog-decisions-test/pkg/closer"
)

type App struct {
	Config  *config.AppConfig
	Logger  *slog.Logger
	Closer  *closer.Closer
	Handler *handler.App
}
