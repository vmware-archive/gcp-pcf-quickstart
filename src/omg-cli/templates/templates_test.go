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

func fixturesDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "fixtures")
}

func readFixture(f string) []byte {
	in, err := ioutil.ReadFile(filepath.Join(fixturesDir(), f))
	Expect(err).ToNot(HaveOccurred())

	return in
}

func readYAML(f string) map[string]interface{} {
	out := make(map[string]interface{})
	err := yaml.Unmarshal(readFixture(f), out)
	Expect(err).ToNot(HaveOccurred())

	return out
}

func directorMatchesFixture(director ompattern.Director, suffix string) {
	template, err := director.ToTemplate().Evaluate(true)
	Expect(err).ToNot(HaveOccurred())
	Expect(template).To(MatchYAML(readFixture(fmt.Sprintf("bosh/%s.yml", suffix))))
}

func tilesMatchFixtures(tiles []ompattern.Tile, suffix string) {
	for _, tile := range tiles {
		template, err := tile.ToTemplate().Evaluate(true)
		Expect(err).ToNot(HaveOccurred())
		Expect(template).To(MatchYAML(readFixture(fmt.Sprintf("%s/%s.yml", tile.Name, suffix))))
	}
}

var _ = Describe("GetPattern", func() {
	var (
		pattern        ompattern.Pattern
		healthwatch    bool
		smallFootPrint bool
		varsFile       string
	)
	JustBeforeEach(func() {
		var err error
		pattern, err = GetPattern(&config.EnvConfig{
			SmallFootprint:     smallFootPrint,
			IncludeHealthwatch: healthwatch,
		}, readYAML(varsFile), true)
		Expect(err).ToNot(HaveOccurred())
		err = pattern.Validate(true)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when small-footprint is enabled", func() {
		BeforeEach(func() {
			smallFootPrint = true
			healthwatch = false
			varsFile = "vars-small.yml"
		})
		It("renders tile configs", func() {
			directorMatchesFixture(pattern.Director, "small")
			tilesMatchFixtures(pattern.Tiles, "small")
		})
	})

	Context("when small-footprint and healthwatch are enabled", func() {
		BeforeEach(func() {
			smallFootPrint = true
			healthwatch = true
			varsFile = "vars-small.yml"
		})
		It("renders tile configs", func() {
			directorMatchesFixture(pattern.Director, "small-healthwatch")
			tilesMatchFixtures(pattern.Tiles, "small")
		})
	})

	Context("when small-footprint is disabled", func() {
		BeforeEach(func() {
			smallFootPrint = false
			healthwatch = false
			varsFile = "vars.yml"
		})
		It("renders tile configs", func() {
			directorMatchesFixture(pattern.Director, "full")
			tilesMatchFixtures(pattern.Tiles, "full")
		})
	})
	Context("when small-footprint is disabled and healthwatch enabled", func() {
		BeforeEach(func() {
			smallFootPrint = false
			healthwatch = true
			varsFile = "vars.yml"
		})
		It("renders tile configs", func() {
			directorMatchesFixture(pattern.Director, "full-healthwatch")
			tilesMatchFixtures(pattern.Tiles, "full")
		})
	})
})
