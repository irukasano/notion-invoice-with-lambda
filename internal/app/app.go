package app

import (
	"context"

	"github.com/irukasano/notion-invoice-with-lambda/internal/logger"
)

type App struct {
	logger *logger.Logger
}

type Runner interface {
	Run(ctx context.Context) error
	Logger() *logger.Logger
}

func New(log *logger.Logger) *App {
	return &App{logger: log}
}

func (a *App) Run(ctx context.Context) error {
	// TODO: implement business logic
	_ = ctx
	return nil
}

func (a *App) Logger() *logger.Logger {
	return a.logger
}
