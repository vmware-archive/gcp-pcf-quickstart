package pivnet

import (
	"fmt"
	"io/ioutil"
	"os"

	"omg-cli/config"
)

type TileCache struct {
	Dir string
}

func (tc *TileCache) FileName(tile config.PivnetMetadata) string {
	return fmt.Sprintf("%s-%d-%d.pivotal", tile.Name, tile.ReleaseId, tile.FileId)
}

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
			return os.OpenFile(file.Name(), os.O_RDONLY, 0)
		}
	}

	return nil, nil
}
