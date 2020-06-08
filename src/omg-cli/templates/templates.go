package templates

//go:generate go run -mod=vendor -tags=dev generate.go

import (
	"omg-cli/config"

	"github.com/starkandwayne/om-tiler/pattern"
)

func GetPattern(envCfg *config.EnvConfig, vars map[string]interface{}, varsStore string, expectAllKeys bool) (pattern.Pattern, error) {
	var opsFiles []string
	if !envCfg.SmallFootprint {
		opsFiles = append(opsFiles, "options/full.yml")
	}
	if envCfg.IncludeHealthwatch {
		opsFiles = append(opsFiles, "options/healthwatch.yml")
	}
	if envCfg.IncludeHealthwatch && !envCfg.SmallFootprint {
		opsFiles = append(opsFiles, "options/healthwatch-full.yml")
	}
	return pattern.NewPattern(pattern.Template{
		Store:    Templates,
		Manifest: "deployment.yml",
		OpsFiles: opsFiles,
		Vars:     vars,
	}, varsStore, expectAllKeys)
}
