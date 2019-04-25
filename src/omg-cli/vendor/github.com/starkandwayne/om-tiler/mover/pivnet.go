package mover

import (
	"context"
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

// PivnetClient interface to be implemented by packages which perform pivnet downloads
//go:generate counterfeiter . PivnetClient
type PivnetClient interface {
	DownloadFile(context.Context, pattern.PivnetFile, string) (*os.File, error)
}
