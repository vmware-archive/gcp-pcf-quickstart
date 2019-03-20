package templates

//go:generate go run -tags=dev generate.go

import (
	"omg-cli/config"

	"github.com/starkandwayne/om-tiler/pattern"
)

func GetPattern(envCfg *config.EnvConfig, vars map[string]interface{}) (pattern.Pattern, error) {
	var opsFiles []string
	if envCfg.SmallFootprint {
		opsFiles = append(opsFiles, "options/small-footprint.yml")
	}
	if envCfg.IncludeHealthwatch {
		opsFiles = append(opsFiles, "options/healthwatch.yml")
	}
	return pattern.NewPattern(pattern.Template{
		Store:    Templates,
		Manifest: "deployment.yml",
		OpsFiles: opsFiles,
		Vars:     vars,
	})
}
