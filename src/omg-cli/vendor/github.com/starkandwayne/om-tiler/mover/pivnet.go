package mover

import (
	"context"
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	DownloadFile(context.Context, pattern.PivnetFile, string) (*os.File, error)
}
