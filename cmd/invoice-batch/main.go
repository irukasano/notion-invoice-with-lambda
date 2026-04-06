package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/irukasano/notion-invoice-with-lambda/internal/app"
	"github.com/irukasano/notion-invoice-with-lambda/internal/logger"
)

type handler struct {
	app app.Runner
}

func main() {
	log := logger.New()
	h := handler{
		app: app.New(log),
	}
	lambda.Start(h.handle)
}

func (h handler) handle(ctx context.Context) (any, error) {
	entry := h.app.Logger().WithContext(ctx)
	entry.Info("invoice-batch started", nil)

	if err := h.app.Run(ctx); err != nil {
		entry.Error("invoice-batch failed", map[string]any{"error": err.Error()})
		return nil, err
	}

	entry.Info("invoice-batch finished", nil)
	return nil, nil
}
