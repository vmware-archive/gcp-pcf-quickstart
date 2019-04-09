package tiler

import (
	"context"

	goflow "github.com/kamildrazkiewicz/go-flow"
	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/steps"
)

func (c *Tiler) Build(ctx context.Context, p pattern.Pattern, skipApplyChanges bool) error {
	if err := p.Validate(true); err != nil {
		return err
	}

	pollTillOnline := steps.Step(ctx, "pollTillOnline", c.client.PollTillOnline)
	configureAuthentication := steps.Step(ctx, "configureAuthentication", c.client.ConfigureAuthentication)
	configureDirector := steps.Step(ctx, "configureDirector", func(ctx context.Context) error {
		err := c.configureDirector(ctx, p.Director)
		if err != nil {
			return err
		}

		if !skipApplyChanges {
			return c.client.ApplyChanges(ctx)
		}
		return nil
	})

	ensureFilesUploaded := steps.Step(ctx, "ensureFilesUploaded", func(ctx context.Context) error {
		for _, tile := range p.Tiles {
			if err := c.ensureFilesUploaded(ctx, tile); err != nil {
				return err
			}
		}
		return nil
	})

	configureTiles := steps.Step(ctx, "configureTiles", func(ctx context.Context) error {
		for _, tile := range p.Tiles {
			err := c.client.StageProduct(ctx, tile)
			if err != nil {
				return err
			}

			err = c.configureProduct(ctx, tile)
			if err != nil {
				return err
			}
		}

		if !skipApplyChanges {
			return c.client.ApplyChanges(ctx)
		}
		return nil
	})

	_, err := goflow.New().
		Add("pollTillOnline", nil, pollTillOnline).
		Add("configureAuthentication", []string{"pollTillOnline"}, configureAuthentication).
		Add("configureDirector", []string{"configureAuthentication"}, configureDirector).
		Add("ensureFilesUploaded", []string{"configureAuthentication"}, ensureFilesUploaded).
		Add("configureTiles", []string{"configureDirector", "ensureFilesUploaded"}, configureTiles).
		Do()

	return err
}

func (c *Tiler) ensureFilesUploaded(ctx context.Context, t pattern.Tile) error {
	ok, err := c.client.FilesUploaded(ctx, t)
	if err != nil {
		return err
	}
	if ok {
		c.logger(ctx).Printf("files for %s/%s already uploaded skipping download",
			t.Name, t.Version)
		return nil
	}

	product, err := c.mover.Get(ctx, t.Product)
	if err != nil {
		return err
	}

	if err = c.client.UploadProduct(ctx, product); err != nil {
		return err
	}

	stemcell, err := c.mover.Get(ctx, t.Stemcell)
	if err != nil {
		return err
	}

	return c.client.UploadStemcell(ctx, stemcell)
}

func (c *Tiler) configureProduct(ctx context.Context, t pattern.Tile) error {
	tpl, err := t.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return c.client.ConfigureProduct(ctx, tpl)
}

func (c *Tiler) configureDirector(ctx context.Context, d pattern.Director) error {
	tpl, err := d.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return c.client.ConfigureDirector(ctx, tpl)
}
