package tiler

import (
	"fmt"
	"path/filepath"

	goflow "github.com/kamildrazkiewicz/go-flow"
	"github.com/starkandwayne/om-tiler/pattern"
)

func (c *Tiler) Build(p pattern.Pattern, skipApplyChanges bool) error {
	if err := p.Validate(true); err != nil {
		return err
	}

	pollTillOnline := func(r map[string]interface{}) (interface{}, error) {
		err := c.client.PollTillOnline()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	configureAuthentication := func(r map[string]interface{}) (interface{}, error) {
		err := c.client.ConfigureAuthentication()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	configureDirector := func(r map[string]interface{}) (interface{}, error) {
		err := c.configureDirector(p.Director)
		if err != nil {
			return nil, err
		}

		if !skipApplyChanges {
			err = c.client.ApplyChanges()
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	ensureFilesUploaded := func(r map[string]interface{}) (interface{}, error) {
		for _, tile := range p.Tiles {
			if err := c.ensureFilesUploaded(tile); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	configureTiles := func(r map[string]interface{}) (interface{}, error) {
		for _, tile := range p.Tiles {
			err := c.client.StageProduct(tile)
			if err != nil {
				return nil, err
			}

			err = c.configureProduct(tile)
			if err != nil {
				return nil, err
			}
		}

		if !skipApplyChanges {
			err := c.client.ApplyChanges()
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	_, err := goflow.New().
		Add("pollTillOnline", nil, pollTillOnline).
		Add("configureAuthentication", []string{"pollTillOnline"}, configureAuthentication).
		Add("configureDirector", []string{"configureAuthentication"}, configureDirector).
		Add("ensureFilesUploaded", []string{"configureAuthentication"}, ensureFilesUploaded).
		Add("configureTiles", []string{"configureDirector", "ensureFilesUploaded"}, configureTiles).
		Do()

	return err
}

func (c *Tiler) ensureFilesUploaded(t pattern.Tile) error {
	ok, err := c.client.FilesUploaded(t)
	if err != nil {
		return err
	}
	if ok {
		c.logger.Printf("files for %s/%s already uploaded skipping download",
			t.Name, t.Version)
		return nil
	}

	product, err := c.mover.Get(t.Product)
	if err != nil {
		return err
	}

	if err = c.client.UploadProduct(product); err != nil {
		return err
	}

	stemcell, err := c.mover.Get(t.Stemcell)
	if err != nil {
		return err
	}

	return c.client.UploadStemcell(stemcell)
}

func (c *Tiler) configureProduct(t pattern.Tile) error {
	tpl, err := t.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return c.client.ConfigureProduct(tpl)
}

func (c *Tiler) configureDirector(d pattern.Director) error {
	tpl, err := d.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return c.client.ConfigureDirector(tpl)
}

func findFileInDir(dir, glob string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, glob))
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", fmt.Errorf("no file found for %s in %s", glob, dir)
	}
	return files[0], nil
}
