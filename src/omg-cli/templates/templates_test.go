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

func fixturePath(f string) string {
	return filepath.Join(fixturesDir(), f)
}

func readFixture(f string) []byte {
	in, err := ioutil.ReadFile(fixturePath(f))
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
		pattern        ompattern.Pattern
		healthwatch    bool
		smallFootPrint bool
		varsFile       string
		varsStore      string
	)
	JustBeforeEach(func() {
		var err error
		pattern, err = GetPattern(&config.EnvConfig{
			SmallFootprint:     smallFootPrint,
			IncludeHealthwatch: healthwatch,
		}, readYAML(varsFile), fixturePath(varsStore), true)
		Expect(err).ToNot(HaveOccurred())
		err = pattern.Validate(true)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when small-footprint is enabled", func() {
		BeforeEach(func() {
			smallFootPrint = true
			healthwatch = false
			varsFile = "vars-small.yml"
			varsStore = "creds.yml"
		})
		It("renders tile configs", func() {
			pattern.MatchesFixtures(ompattern.Fixtures{
				Dir:            fixturesDir(),
				DirectorSuffix: "small",
				TilesSuffix:    "small",
			})
		})
	})

	Context("when small-footprint and healthwatch are enabled", func() {
		BeforeEach(func() {
			smallFootPrint = true
			healthwatch = true
			varsFile = "vars-small.yml"
			varsStore = "creds.yml"
		})
		It("renders tile configs", func() {
			pattern.MatchesFixtures(ompattern.Fixtures{
				Dir:            fixturesDir(),
				DirectorSuffix: "small-healthwatch",
				TilesSuffix:    "small",
			})
		})
	})

	Context("when small-footprint is disabled", func() {
		BeforeEach(func() {
			smallFootPrint = false
			healthwatch = false
			varsFile = "vars.yml"
			varsStore = "creds.yml"
		})
		It("renders tile configs", func() {
			pattern.MatchesFixtures(ompattern.Fixtures{
				Dir:            fixturesDir(),
				DirectorSuffix: "full",
				TilesSuffix:    "full",
			})
		})
	})
	Context("when small-footprint is disabled and healthwatch enabled", func() {
		BeforeEach(func() {
			smallFootPrint = false
			healthwatch = true
			varsFile = "vars.yml"
			varsStore = "creds.yml"
		})
		It("renders tile configs", func() {
			pattern.MatchesFixtures(ompattern.Fixtures{
				Dir:            fixturesDir(),
				DirectorSuffix: "full-healthwatch",
				TilesSuffix:    "full",
			})
		})
	})
})
