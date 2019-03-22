package templates_test

import (
	"fmt"
	"io/ioutil"
	"omg-cli/config"
	. "omg-cli/templates"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	ompattern "github.com/starkandwayne/om-tiler/pattern"
	"github.com/thadc23/yamldiff/differ"
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
		pattern         ompattern.Pattern
		healthwatch     bool
		smallfootprint  bool
		varsfile        string
		tileMatchesMock func(string, string) string
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

		tileMatchesMock = func(t string, f string) string {
			var tile ompattern.Tile
			for _, tile = range pattern.Tiles {
				if tile.Name == t {
					break
				}
			}
			if tile.Name != t {
				panic(fmt.Sprintf("expected to find tile: %s", t))
			}
			actualRaw, err := tile.ToTemplate().Evaluate(true)
			Expect(err).ToNot(HaveOccurred())

			actual, err := ioutil.TempFile(os.TempDir(), t)
			Expect(err).ToNot(HaveOccurred())
			defer os.Remove(actual.Name())
			_, err = actual.Write(actualRaw)
			Expect(err).ToNot(HaveOccurred())
			Expect(actual.Close()).ToNot(HaveOccurred())

			diff := differ.NewDiffer(actual.Name(),
				filepath.Join(fixturesDir(), f), false).ComputeDiff()
			if diff != "" {
				fmt.Println(string(actualRaw))
			}
			return diff

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
			Expect(tileMatchesMock("cf", "cf-smallfootprint.yml")).To(Equal(""))
			Expect(tileMatchesMock("stackdriver-nozzle", "stackdriver.yml")).To(Equal(""))
			Expect(tileMatchesMock("gcp-service-broker", "service-broker.yml")).To(Equal(""))
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
			Expect(tileMatchesMock("cf", "cf-smallfootprint.yml")).To(Equal(""))
			Expect(tileMatchesMock("stackdriver-nozzle", "stackdriver.yml")).To(Equal(""))
			Expect(tileMatchesMock("gcp-service-broker", "service-broker.yml")).To(Equal(""))
			Expect(tileMatchesMock("p-healthwatch", "p-healthwatch.yml")).To(Equal(""))
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
			Expect(tileMatchesMock("cf", "cf.yml")).To(Equal(""))
			Expect(tileMatchesMock("stackdriver-nozzle", "stackdriver.yml")).To(Equal(""))
			Expect(tileMatchesMock("gcp-service-broker", "service-broker.yml")).To(Equal(""))
			Expect(tileMatchesMock("p-healthwatch", "p-healthwatch.yml")).To(Equal(""))
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
