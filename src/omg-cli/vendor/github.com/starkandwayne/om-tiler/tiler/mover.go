package tiler

import (
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

//go:generate counterfeiter . Mover
type Mover interface {
	Get(pattern.PivnetFile) (*os.File, error)
	Cache(pattern.PivnetFile) error
}
