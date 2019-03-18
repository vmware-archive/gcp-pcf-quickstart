package mover

import (
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	DownloadFile(pattern.PivnetFile, string) (*os.File, error)
}
