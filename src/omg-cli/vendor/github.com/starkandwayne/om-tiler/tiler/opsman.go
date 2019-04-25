package tiler

import (
	"context"
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

// OpsmanClient responsible for interacting with the OpsManager API
//go:generate counterfeiter . OpsmanClient
type OpsmanClient interface {
	PollTillOnline(context.Context) error
	ConfigureAuthentication(context.Context) error
	UploadProduct(context.Context, *os.File) error
	UploadStemcell(context.Context, *os.File) error
	FilesUploaded(context.Context, pattern.Tile) (bool, error)
	StageProduct(context.Context, pattern.Tile) error
	ConfigureDirector(context.Context, []byte) error
	ConfigureProduct(context.Context, []byte) error
	ApplyChanges(context.Context) error
	DeleteInstallation(context.Context) error
}
