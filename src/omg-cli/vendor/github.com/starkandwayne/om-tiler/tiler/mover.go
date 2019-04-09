package tiler

import (
	"context"
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

//go:generate counterfeiter . Mover
type Mover interface {
	Get(context.Context, pattern.PivnetFile) (*os.File, error)
	Cache(context.Context, pattern.PivnetFile) error
}
