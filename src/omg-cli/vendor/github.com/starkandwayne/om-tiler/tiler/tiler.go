package tiler

import (
	"fmt"
	"log"
)

type Tiler struct {
	client OpsmanClient
	logger *log.Logger
	mover  Mover
}

func NewTiler(c OpsmanClient, m Mover, l *log.Logger) (*Tiler, error) {
	l.SetPrefix(fmt.Sprintf("%s[OM Tiler] ", l.Prefix()))
	return &Tiler{client: c, mover: m, logger: l}, nil
}
