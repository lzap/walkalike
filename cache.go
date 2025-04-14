package walkalike

import (
	"fmt"
	"os"
	"path/filepath"
)

type Cache struct {
	root string
}

func NewCache(root string) *Cache {
	return &Cache{
		root: root,
	}
}

func (c *Cache) Ensure() error {
	if _, err := os.Stat(c.root); os.IsNotExist(err) {
		if err := os.MkdirAll(c.root, 0750); err != nil {
			return fmt.Errorf("failed to create cache directory: %w", err)
		}
	}

	return nil
}

func (c *Cache) Get(path string, info os.FileInfo) (*Index, error) {
	cacheFile, err := ChecksumCache(path, info.Size(), info.ModTime())
	if err != nil {
		return nil, err
	}
	cachePath := filepath.Join(c.root, cacheFile)

	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		ix := &Index{}
		f, err := os.Open(cachePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		err = ix.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("error decoding cache index: %w", err)
		}

		return ix, nil
	}

	return nil, nil
}

func (c *Cache) Put(path string, info os.FileInfo, index *Index) error {
	cacheFile, err := ChecksumCache(path, info.Size(), info.ModTime())
	if err != nil {
		return err
	}
	cachePath := filepath.Join(c.root, cacheFile)

	f, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = index.Encode(f, filepath.Base(path))
	if err != nil {
		return fmt.Errorf("error encoding cache index: %w", err)
	}

	return nil
}
