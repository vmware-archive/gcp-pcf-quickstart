package tiler

import (
	"context"
	"fmt"

	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/steps"
)

const (
	// StepWaitOpsmanOnline invoked [Build, Delete] when OpsMan API is available
	StepWaitOpsmanOnline string = "WaitOpsmanOnline"
	// StepConfigureAuthentication invoked [Build, Delete] when Authentication has been configured
	StepConfigureAuthentication = "ConfigureAuthentication"
	// StepConfigureDirector invoked [Build] when Director has been configured
	StepConfigureDirector = "ConfigureDirector"
	// StepDeployDirector invoked [Build] when Director has been deployed
	StepDeployDirector = "DeployDirector"
	// StepUploadFiles invoked [Build] when all PivnetFiles have been uploaded
	StepUploadFiles = "UploadFiles"
	// StepConfigureTiles invoked [Build] when all Tiles have ben configured
	StepConfigureTiles = "ConfigureTiles"
	// StepApplyChanges invoked [Build] when all Tiles have been deployed
	StepApplyChanges = "ApplyChanges"
	// StepDeleteInstallation invoked [Delete] when all Tiles have been deleted
	StepDeleteInstallation = "DeleteInstallation"
	retry                  = 5
)

func stepUploadFilesName(tile pattern.Tile) string {
	return fmt.Sprintf(
		"UploadFiles(%s/%s)", tile.Name, tile.Version)
}

func stepConfigureTileName(tile pattern.Tile) string {
	return fmt.Sprintf(
		"ConfigureTiles(%s/%s)", tile.Name, tile.Version)
}

func (t *Tiler) stepPollTillOnline() steps.Step {
	return steps.Step{
		Name: StepWaitOpsmanOnline,
		Do:   t.client.PollTillOnline,
	}
}

func (t *Tiler) stepConfigureAuthentication() steps.Step {
	return steps.Step{
		Name:      StepConfigureAuthentication,
		DependsOn: []string{StepWaitOpsmanOnline},
		Do:        t.client.ConfigureAuthentication,
		Retry:     retry,
	}
}

func (t *Tiler) stepConfigureDirector(d pattern.Director) steps.Step {
	return steps.Step{
		Name:      StepConfigureDirector,
		DependsOn: []string{StepConfigureAuthentication},
		Do: func(ctx context.Context) error {
			return t.doConfigureDirector(ctx, d)
		},
		Retry: retry,
	}
}

func (t *Tiler) stepDeployDirector(skipApplyChanges bool) steps.Step {
	s := t.stepApplyChanges(skipApplyChanges)
	s.Name = StepDeployDirector
	s.DependsOn = []string{StepConfigureDirector}
	return s
}

func (t *Tiler) stepUploadFiles(tiles []pattern.Tile) (out []steps.Step) {
	var dependsOn []string
	for _, tile := range tiles {
		tile := tile
		dependsOn = append(dependsOn, stepUploadFilesName(tile))
		out = append(out, steps.Step{
			Name:      stepUploadFilesName(tile),
			DependsOn: []string{StepConfigureAuthentication},
			Do: func(ctx context.Context) error {
				return t.doUploadFiles(ctx, tile)
			},
			Retry: retry,
		})
	}
	return append(out, steps.Step{
		Name:      StepUploadFiles,
		DependsOn: dependsOn,
	})
}

func (t *Tiler) stepConfigureTiles(tiles []pattern.Tile) (out []steps.Step) {
	var dependsOn []string
	for _, tile := range tiles {
		tile := tile
		dependsOn = append(dependsOn, stepConfigureTileName(tile))
		out = append(out, steps.Step{
			Name:      stepConfigureTileName(tile),
			DependsOn: []string{StepDeployDirector, stepUploadFilesName(tile)},
			Do: func(ctx context.Context) error {
				return t.doConfigureTile(ctx, tile)
			},
			Retry: retry,
		})
	}
	return append(out, steps.Step{
		Name:      StepConfigureTiles,
		DependsOn: dependsOn,
	})
}

func (t *Tiler) stepApplyChanges(skipApplyChanges bool) steps.Step {
	s := steps.Step{
		Name:      StepApplyChanges,
		DependsOn: []string{StepConfigureTiles},
	}
	if !skipApplyChanges {
		s.Do = t.client.ApplyChanges
	}
	return s

}

func (t *Tiler) stepDeleteInstallation() steps.Step {
	return steps.Step{
		Name:      StepApplyChanges,
		DependsOn: []string{StepConfigureAuthentication},
		Do:        t.client.DeleteInstallation,
	}
}

func (t *Tiler) doConfigureDirector(ctx context.Context, d pattern.Director) error {
	tpl, err := d.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return t.client.ConfigureDirector(ctx, tpl)
}

func (t *Tiler) doUploadFiles(ctx context.Context, tile pattern.Tile) error {
	ok, err := t.client.FilesUploaded(ctx, tile)
	if err != nil {
		return err
	}
	if ok {
		t.logger(ctx).Printf("files for %s/%s already uploaded skipping download",
			tile.Name, tile.Version)
		return nil
	}

	product, err := t.mover.Get(ctx, tile.Product)
	if err != nil {
		return err
	}

	if err = t.client.UploadProduct(ctx, product); err != nil {
		return err
	}

	stemcell, err := t.mover.Get(ctx, tile.Stemcell)
	if err != nil {
		return err
	}

	return t.client.UploadStemcell(ctx, stemcell)
}

func (t *Tiler) doConfigureTile(ctx context.Context, tile pattern.Tile) error {
	err := t.client.StageProduct(ctx, tile)
	if err != nil {
		return err
	}

	tpl, err := tile.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}
	return t.client.ConfigureProduct(ctx, tpl)
}
