package app

import (
	"context"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/nut_parser/internal/http/server"
	metricsNut "github.com/andreyAKor/nut_parser/internal/metrics/nut"
)

var _ io.Closer = (*App)(nil)

type App struct {
	srv        *server.Server
	nutMetrics *metricsNut.Metric
}

func New(srv *server.Server, nutMetrics *metricsNut.Metric) (*App, error) {
	return &App{
		srv:        srv,
		nutMetrics: nutMetrics,
	}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.srv.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("http-server listen fail")
		}
	}()
	go func() {
		if err := a.nutMetrics.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("nut metrics running fail")
		}
	}()

	return nil
}

// Close application.
func (a *App) Close() error {
	return nil
}
