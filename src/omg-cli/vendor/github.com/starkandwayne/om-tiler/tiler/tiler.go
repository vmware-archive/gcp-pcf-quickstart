package tiler

import (
	"log"
)

type Tiler struct {
	client OpsmanClient
	logger *log.Logger
	mover  Mover
}

func NewTiler(c OpsmanClient, m Mover, l *log.Logger) (*Tiler, error) {
	return &Tiler{client: c, mover: m, logger: l}, nil
}
