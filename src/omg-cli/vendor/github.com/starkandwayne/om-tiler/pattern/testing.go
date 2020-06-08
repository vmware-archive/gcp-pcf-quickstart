package pattern

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/onsi/gomega"
)

type Fixtures struct {
	Dir            string
	DirectorSuffix string
	TilesSuffix    string
}

func (p *Pattern) MatchesFixtures(f Fixtures) {
	template, err := p.Director.ToTemplate().Evaluate(true)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(template).To(gomega.MatchYAML(f.readFixture(fmt.Sprintf("director/%s.yml", f.DirectorSuffix))))
	for _, tile := range p.Tiles {
		template, err := tile.ToTemplate().Evaluate(true)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		fixturePath := fmt.Sprintf("%s/%s.yml", tile.Name, f.TilesSuffix)
		if os.Getenv("UPDATE_FIXTURES") != "" {
			f.writeFixture(fixturePath, template)
		}
		gomega.Expect(template).To(gomega.MatchYAML(f.readFixture(fixturePath)))
	}
}

func (f *Fixtures) readFixture(name string) []byte {
	in, err := ioutil.ReadFile(filepath.Join(f.Dir, name))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	return in
}

func (f *Fixtures) writeFixture(name string, data []byte) {
	err := ioutil.WriteFile(filepath.Join(f.Dir, name), data, 0644)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
}
