package templates_test

import (
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

var _ = Describe("GetPattern", func() {
	var (
		pattern            ompattern.Pattern
		healthwatch        bool
		smallfootprint     bool
		varsfile           string
		tileMatchesFixture func(string, string)
	)
	JustBeforeEach(func() {
		var err error
		pattern, err = GetPattern(&config.EnvConfig{
			SmallFootprint:     smallfootprint,
			IncludeHealthwatch: healthwatch,
		}, readYAML(varsfile), true)
		Expect(err).ToNot(HaveOccurred())
		err = pattern.Validate(true)
		Expect(err).ToNot(HaveOccurred())

		tileMatchesFixture = func(name, fixture string) {
			var tile ompattern.Tile
			for _, tile = range pattern.Tiles {
				if tile.Name == name {
					break
				}
			}
			Expect(tile.Name).ToNot(Equal(""), "Expected to find tile")
			result, err := tile.ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(MatchYAML(readFixture(fixture)))
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
			Expect(director).To(MatchYAML(readFixture("bosh/smallfootprint.yml")))
			tileMatchesFixture("cf", "cf/smallfootprint.yml")
			tileMatchesFixture("stackdriver-nozzle", "stackdriver/smallfootprint.yml")
			tileMatchesFixture("gcp-service-broker", "service-broker/smallfootprint.yml")
		})
	})

	Context("when small-footprint and healthwatch are enabled", func() {
		BeforeEach(func() {
			smallfootprint = true
			healthwatch = true
			varsfile = "vars-smallfootprint.yml"
		})
		It("renders tile configs", func() {
			director, err := pattern.Director.ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())
			Expect(director).To(MatchYAML(readFixture("bosh/smallfootprint-healthwatch.yml")))
			tileMatchesFixture("cf", "cf/smallfootprint.yml")
			tileMatchesFixture("stackdriver-nozzle", "stackdriver/smallfootprint.yml")
			tileMatchesFixture("gcp-service-broker", "service-broker/smallfootprint.yml")
			tileMatchesFixture("p-healthwatch", "healthwatch/smallfootprint.yml")
		})
	})

	Context("when small-footprint is disabled", func() {
		BeforeEach(func() {
			smallfootprint = false
			healthwatch = false
			varsfile = "vars.yml"
		})
		It("renders tile configs", func() {
			director, err := pattern.Director.ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())
			Expect(director).To(MatchYAML(readFixture("bosh/full.yml")))
			tileMatchesFixture("cf", "cf/full.yml")
			tileMatchesFixture("stackdriver-nozzle", "stackdriver/full.yml")
			tileMatchesFixture("gcp-service-broker", "service-broker/full.yml")
		})
	})
	Context("when small-footprint is disabled and healthwatch enabled", func() {
		BeforeEach(func() {
			smallfootprint = false
			healthwatch = true
			varsfile = "vars.yml"
		})
		It("renders tile configs", func() {
			director, err := pattern.Director.ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())
			Expect(director).To(MatchYAML(readFixture("bosh/full-healthwatch.yml")))
			tileMatchesFixture("cf", "cf/full.yml")
			tileMatchesFixture("stackdriver-nozzle", "stackdriver/full.yml")
			tileMatchesFixture("gcp-service-broker", "service-broker/full.yml")
			tileMatchesFixture("p-healthwatch", "healthwatch/full.yml")
		})
	})
})
