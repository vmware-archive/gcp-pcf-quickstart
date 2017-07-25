package commands

import (
	"log"

	"omg-cli/config"
	"omg-cli/omg/tiles"
	"omg-cli/omg/tiles/ert"
	"omg-cli/omg/tiles/gcp_director"
	"omg-cli/omg/tiles/service_broker"
	"omg-cli/omg/tiles/stackdriver_nozzle"

	"github.com/alecthomas/kingpin"
)

const (
	defaultUsername         = "foo"
	defaultPassword         = "foobar"
	defaultDecryptionPhrase = "foobar"
	defaultSkipSSLVerify    = "true"
)

var selectedTiles = []tiles.TileInstaller{
	gcp_director.Tile{},
	ert.Tile{},
	stackdriver_nozzle.Tile{},
	service_broker.Tile{},
}

type register interface {
	register(app *kingpin.Application)
}

func Configure(logger *log.Logger, app *kingpin.Application) {
	cmds := []register{
		&PushTilesCommand{logger: logger},
		&Deploy{logger: logger},
		&BootstrapDeployCommand{logger: logger},
		&BootstrapPushTilesCommand{logger: logger},
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

func registerOpsManagerFlags(c *kingpin.CmdClause, cfg *config.OpsManagerCredentials) {
	c.Flag("opsman-username", "Username for Ops Manager").Default(defaultUsername).StringVar(&cfg.Username)
	c.Flag("opsman-password", "Password for Ops Manager").Default(defaultPassword).StringVar(&cfg.Password)
	c.Flag("opsman-decryption-phrase", "Decryption Phrase for Ops Manager").Default(defaultDecryptionPhrase).StringVar(&cfg.DecryptionPhrase)
	c.Flag("opsman-skip-ssl-verification", "Skip SSL Validation for Ops Manager").Default(defaultSkipSSLVerify).BoolVar(&cfg.SkipSSLVerification)
}

func registerTerraformConfigFlag(c *kingpin.CmdClause, path *string) {
	c.Flag("terraform-output-path", "JSON output from terraform state for deployment").Default("env.json").StringVar(path)
}

func registerPivnetFlag(c *kingpin.CmdClause, apiToken *string) {
	c.Flag("pivnet-api-token", "Look for 'API TOKEN' at https://network.pivotal.io/users/dashboard/edit-profile.").Required().StringVar(apiToken)
}
