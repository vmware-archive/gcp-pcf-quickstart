package pattern

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

// Pattern can be given to tiler to lay tiles
type Pattern struct {
	Director  Director      `yaml:"director" validate:"required,dive"`
	Tiles     []Tile        `yaml:"tiles" validate:"required,dive"`
	Variables []interface{} `yaml:"variables"`
}

// NewPattern will render a Template using optionally a given varsStore
func NewPattern(t Template, varsStore string, expectAllKeys bool) (p Pattern, err error) {
	t.VarsStore = varsStore
	db, err := t.Evaluate(expectAllKeys)
	if err != nil {
		return Pattern{}, fmt.Errorf("could not render pattern: %v", err)
	}

	if err = yaml.UnmarshalStrict(db, &p); err != nil {
		return Pattern{}, fmt.Errorf("could not unmarshal rendered pattern: %v", err)
	}

	if p.Director.Vars == nil {
		p.Director.Vars = make(map[string]interface{})
	}
	mergeVars(p.Director.Vars, t.Vars)
	p.Director.Store = t.Store

	for i := range p.Tiles {
		if p.Tiles[i].Vars == nil {
			p.Tiles[i].Vars = make(map[string]interface{})
		}
		mergeVars(p.Tiles[i].Vars, t.Vars)
		p.Tiles[i].Store = t.Store
	}

	return p, err
}

// Validate check if all required fields are provided
func (p *Pattern) Validate(expectAllKeys bool) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	err := validate.Struct(p)
	if err != nil {
		return fmt.Errorf("pattern.Pattern has error(s):\n%+v\n", err)
	}

	_, err = p.Director.ToTemplate().Evaluate(expectAllKeys)
	if err != nil {
		return fmt.Errorf("Director interpolation error(s):\n%+v\n", err)
	}

	for _, tile := range p.Tiles {
		_, err = tile.ToTemplate().Evaluate(expectAllKeys)
		if err != nil {
			return fmt.Errorf("Tile %s interpolation error(s):\n%+v\n", tile.Name, err)
		}
	}

	return nil
}

// Template can be rendered using the bosh interpolate library
type Template struct {
	Manifest  string                 `yaml:"manifest"`
	OpsFiles  []string               `yaml:"ops_files"`
	VarsFiles []string               `yaml:"vars_files"`
	Vars      map[string]interface{} `yaml:"vars"`
	VarsStore string
	Store     http.FileSystem
}

// Director describes how to configure the OpsManager director
type Director Template

// ToTemplate converts a Director to a Template
func (d *Director) ToTemplate() *Template {
	return &Template{
		Manifest:  d.Manifest,
		OpsFiles:  d.OpsFiles,
		VarsFiles: d.VarsFiles,
		Vars:      d.Vars,
		Store:     d.Store,
	}
}

// Tile describes how to configure a tile, as well as where to download it from
type Tile struct {
	Name     string     `yaml:"name" validate:"required"`
	Version  string     `yaml:"version" validate:"required"`
	Product  PivnetFile `yaml:"product" validate:"required,dive"`
	Stemcell PivnetFile `yaml:"stemcell" validate:"required,dive"`
	Template `yaml:",inline"`
}

// ToTemplate converts a Tile to a Template
func (t *Tile) ToTemplate() *Template {
	return &Template{
		Manifest:  t.Manifest,
		OpsFiles:  t.OpsFiles,
		VarsFiles: t.VarsFiles,
		Vars:      t.Vars,
		Store:     t.Store,
	}
}

// PivnetFile references a file on the pivotal network
// or optionally via an URL (when using a caching proxy)
type PivnetFile struct {
	Slug    string `yaml:"product_slug" validate:"required"`
	Version string `yaml:"release_version" validate:"required"`
	Glob    string `yaml:"file_glob" validate:"required"`
	URL     string `yaml:"download_url"`
}

func mergeVars(target map[string]interface{}, source map[string]interface{}) {
	for k, v := range source {
		if _, ok := target[k]; !ok {
			target[k] = v
		}
	}
}
