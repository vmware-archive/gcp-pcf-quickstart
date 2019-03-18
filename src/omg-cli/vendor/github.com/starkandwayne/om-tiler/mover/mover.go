package mover

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/starkandwayne/om-tiler/pattern"
)

type Mover struct {
	client PivnetClient
	logger *log.Logger
	cache  string
}

func NewMover(c PivnetClient, cache string, l *log.Logger) (*Mover, error) {
	l.SetPrefix(fmt.Sprintf("%s[OM Tile Mover] ", l.Prefix()))
	if cache == "" {
		cache = os.TempDir()
	} else {
		err := os.MkdirAll(cache, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &Mover{client: c, cache: cache, logger: l}, nil
}

func (m *Mover) Get(f pattern.PivnetFile) (*os.File, error) {
	file := m.cachedFilePath(f)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if err = m.Cache(f); err != nil {
			return nil, err
		}
	}
	return os.Open(file)
}

func (m *Mover) Cache(f pattern.PivnetFile) error {
	_, err := m.client.DownloadFile(f, m.cachedFilePath(f))
	return err
}

func (m *Mover) cachedFilePath(f pattern.PivnetFile) string {
	return filepath.Join(m.cache, uuid.NewSHA1(uuid.Nil,
		[]byte(strings.Join([]string{f.Slug, f.Version, f.Glob}, "-")),
	).String())
}
