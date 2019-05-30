package mover

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/steps"
)

// Mover downloads and caches pivnet products
type Mover struct {
	client PivnetClient
	logger func(context.Context) *log.Logger
	cache  string
}

// NewMover create a Mover when given a PivnetClient and cache location
func NewMover(c PivnetClient, cache string, l *log.Logger) (*Mover, error) {
	if cache == "" {
		cache = os.TempDir()
	} else {
		err := os.MkdirAll(cache, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("could not create cache dir %s: %v", cache, err)
		}
	}
	log := func(ctx context.Context) *log.Logger {
		return steps.ContextLogger(ctx, l, "[Cache]")
	}
	return &Mover{client: c, cache: cache, logger: log}, nil
}

// Get downloads a given PivnetFiles if not found in cache returns a file
func (m *Mover) Get(ctx context.Context, f pattern.PivnetFile) (*os.File, error) {
	logger := m.logger(ctx)
	ok, file, err := m.cachedFile(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("could not look up cache file: %v", err)
	}
	if ok {
		logger.Printf("using file: %s from cache", file.Name())
		return file, nil
	}

	logger.Printf("file: %s/%s not found in cache", f.Slug, f.Version)
	if err = m.Cache(ctx, f); err != nil {
		return nil, fmt.Errorf("could not cache file: %v", err)
	}
	_, file, err = m.cachedFile(ctx, f)

	return file, err
}

// Cache will store a copy of a PivnetFile in a deterministic cache location
func (m *Mover) Cache(ctx context.Context, f pattern.PivnetFile) error {
	m.logger(ctx).Printf("caching file: %s/%s", f.Slug, f.Version)
	dir, err := m.cachedFileDir(f)
	if err != nil {
		return fmt.Errorf("could not look up cache dir: %v", err)
	}
	_, err = m.client.DownloadFile(ctx, f, dir.Name())
	return err
}

func (m *Mover) cachedFile(ctx context.Context, f pattern.PivnetFile) (bool, *os.File, error) {
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
		m.logger(ctx).Printf("cache corrupted for %s/%s, removing dir %s",
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
