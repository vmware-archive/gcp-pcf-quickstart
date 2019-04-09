package tiler

import (
	"context"
	"log"

	"github.com/starkandwayne/om-tiler/steps"
)

type Tiler struct {
	client OpsmanClient
	logger func(context.Context) *log.Logger
	mover  Mover
}

func NewTiler(c OpsmanClient, m Mover, l *log.Logger) *Tiler {
	log := func(ctx context.Context) *log.Logger {
		return steps.ContextLogger(ctx, l, "[Tiler]")
	}
	return &Tiler{client: c, mover: m, logger: log}
}
