package commands

import (
	"log"

	"omg-cli/omg/tiles"
	"omg-cli/omg/tiles/ert"

	"github.com/alecthomas/kingpin"
)

const (
	defaultUsername         = "foo"
	defaultPassword         = "foobar"
	defaultDecryptionPhrase = "foobar"
	defaultSkipSSLVerify    = "true"
)

var selectedTiles = []tiles.TileInstaller{
	ert.Tile{},
}

type register interface {
	register(app *kingpin.Application)
}

func Configure(logger *log.Logger, app *kingpin.Application) {
	cmds := []register{
		&BakeImageCommand{logger: logger},
		&ConfigureOpsManagerCommand{logger: logger},
	}

	for _, c := range cmds {
		c.register(app)
	}
}

type step func() error

func runSteps(steps []step) error {
	for _, v := range steps {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}
