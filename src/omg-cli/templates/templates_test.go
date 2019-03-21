package templates_test

import (
	"fmt"
	"io/ioutil"
	"omg-cli/config"
	. "omg-cli/templates"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	ompattern "github.com/starkandwayne/om-tiler/pattern"
)

func mocksDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "mocks")
}

func readMock(f string) []byte {
	in, err := ioutil.ReadFile(filepath.Join(mocksDir(), f))
	Expect(err).ToNot(HaveOccurred())

	return in
}

func getVars(f string) map[string]interface{} {
	out := make(map[string]interface{})
	err := yaml.Unmarshal(readMock(f), out)
	Expect(err).ToNot(HaveOccurred())

	return out
}

var _ = Describe("GetPattern", func() {
	var (
		pattern        ompattern.Pattern
		healthwatch    bool
		smallfootprint bool
		varsfile       string
		getTile        func(string) *ompattern.Tile
	)
	JustBeforeEach(func() {
		var err error
		pattern, err = GetPattern(&config.EnvConfig{
			SmallFootprint:     smallfootprint,
			IncludeHealthwatch: healthwatch,
		}, getVars(varsfile))
		Expect(err).ToNot(HaveOccurred())
		err = pattern.Validate(true)
		Expect(err).ToNot(HaveOccurred())
		getTile = func(n string) *ompattern.Tile {
			for _, tile := range pattern.Tiles {
				if tile.Name == n {
					return &tile
				}
			}
			panic(fmt.Sprintf("expected to find tile: %s", n))
		}
	})
	Context("when small-footprint is enabled", func() {
		BeforeEach(func() {
			smallfootprint = true
			healthwatch = false
			varsfile = "vars-smallfootprint.yml"
		})
		It("renders tile configs", func() {
			director, err := pattern.Director.ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())
			Expect(director).To(MatchYAML(readMock("bosh-smallfootprint.yml")))
			cf, err := getTile("cf").ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())
			Expect(cf).To(MatchYAML(readMock("cf-smallfootprint.yml")))
		})
	})
	Context("when small-footprint and healthwatch are enabled", func() {
		BeforeEach(func() {
			smallfootprint = true
			healthwatch = true
			varsfile = "vars-smallfootprint.yml"
		})
	})
	Context("when small-footprint is disabled", func() {
		BeforeEach(func() {
			smallfootprint = false
			healthwatch = false
			varsfile = "vars.yml"
		})
	})
	Context("when small-footprint is disabled and healthwatch enabled", func() {
		BeforeEach(func() {
			smallfootprint = false
			healthwatch = true
			varsfile = "vars.yml"
		})
	})
})
