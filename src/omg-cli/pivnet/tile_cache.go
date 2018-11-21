package pivnet

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"omg-cli/config"
)

// TileCache is a local directory to cache tiles.
type TileCache struct {
	Dir string
}

// FileName returns the tile format for a tile file.
func (tc *TileCache) FileName(tile config.PivnetMetadata) string {
	return fmt.Sprintf("%s-%d-%d.pivotal", tile.Name, tile.ReleaseID, tile.FileID)
}

// Open returns the file contents of a tile, checking the local cache first.
func (tc *TileCache) Open(tile config.PivnetMetadata) (*os.File, error) {
	if tc == nil || tc.Dir == "" {
		return nil, nil
	}

	needle := tc.FileName(tile)

	files, err := ioutil.ReadDir(tc.Dir)
	if err != nil {
		return nil, fmt.Errorf("opening tile cache directory: %v", err)
	}

	for _, file := range files {
		if file.Name() == needle {
			return os.Open(filepath.Join(tc.Dir, file.Name()))
		}
	}

	return nil, nil
}
