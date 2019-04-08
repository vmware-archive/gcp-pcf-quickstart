package tiler

import (
	"os"

	"github.com/starkandwayne/om-tiler/pattern"
)

//go:generate counterfeiter . OpsmanClient
type OpsmanClient interface {
	PollTillOnline() error
	ConfigureAuthentication() error
	UploadProduct(*os.File) error
	UploadStemcell(*os.File) error
	FilesUploaded(pattern.Tile) (bool, error)
	StageProduct(pattern.Tile) error
	ConfigureDirector([]byte) error
	ConfigureProduct([]byte) error
	ApplyChanges() error
	DeleteInstallation() error
}
