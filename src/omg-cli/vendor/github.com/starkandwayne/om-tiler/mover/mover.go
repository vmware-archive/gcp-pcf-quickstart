package mover

import (
	"fmt"
	"io/ioutil"
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
	ok, file, err := m.cachedFile(f)
	if err != nil {
		return nil, err
	}
	if ok {
		m.logger.Printf("using file: %s from cache", file.Name())
		return file, nil
	}

	m.logger.Printf("file: %s/%s not found in cache", f.Slug, f.Version)
	if err = m.Cache(f); err != nil {
		return nil, err
	}
	_, file, err = m.cachedFile(f)

	return file, err
}

func (m *Mover) Cache(f pattern.PivnetFile) error {
	m.logger.Printf("caching file: %s/%s", f.Slug, f.Version)
	dir, err := m.cachedFileDir(f)
	if err != nil {
		return err
	}
	_, err = m.client.DownloadFile(f, dir.Name())
	return err
}

func (m *Mover) cachedFile(f pattern.PivnetFile) (bool, *os.File, error) {
	dir, err := m.cachedFileDir(f)
	if err != nil {
		return false, nil, err
	}

	files, err := ioutil.ReadDir(dir.Name())
	if err != nil {
		return false, nil, err
	}

	if len(files) == 0 {
		return false, nil, nil
	}

	if len(files) > 1 {
		m.logger.Printf("cache corrupted for %s/%s, removing dir %s",
			f.Slug, f.Version, dir.Name())
		err = os.RemoveAll(dir.Name())
		if err != nil {
			return false, nil, fmt.Errorf(
				"cleaning up corruped cached %s: %s", dir.Name(), err)
		}
		_, err := m.cachedFileDir(f)
		if err != nil {
			return false, nil, err
		}
		return false, nil, nil
	}

	filePath := filepath.Join(dir.Name(), files[0].Name())
	file, err := os.Open(filePath)
	if err != nil {
		return false, nil, fmt.Errorf(
			"opening cached file %s: %s", filePath, err)
	}

	return true, file, nil
}

func (m *Mover) cachedFileDir(f pattern.PivnetFile) (*os.File, error) {
	path := filepath.Join(m.cache, uuid.NewSHA1(uuid.Nil,
		[]byte(strings.Join([]string{f.Slug, f.Version, f.Glob}, "-")),
	).String())
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("creating cache dir %s: %s",
			path, err)
	}
	return os.Open(path)
}
