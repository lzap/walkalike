package walkalike

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Indexer struct {
	// ErrFn is called when an error occurs during indexing. The first argument
	// is the path of the file that caused the error, and the second argument
	// is the error itself. The function can be called from multiple goroutines.
	ErrFn func(string, error)

	root string
	ix   *Index
	wg   *sync.WaitGroup
	q    chan string
}

// NewIndexer creates a new Indexer. The root argument is the root of the
// directory tree to be indexed.
func NewIndexer(root string) *Indexer {
	return &Indexer{
		ErrFn: func(path string, err error) {
			// Default error function does nothing
		},
		root: root,
		ix: &Index{
			Tokens: make([]Token, 0, 1024),
		},
		wg: &sync.WaitGroup{},
		q:  make(chan string, 1024),
	}
}

var ErrChecksumSizeMismatch = errors.New("checksum size mismatch")

func (i *Indexer) processFile(ctx context.Context, path string) {
	// get the file info
	stat, err := os.Lstat(path)
	if err != nil {
		i.ErrFn(path, err)
		return
	}

	// calculate path checksum
	pathChecksum := ChecksumPath(path)

	// resolve symlinks
	if stat.Mode()&os.ModeSymlink != 0 {
		dst, err := os.Readlink(path)
		if err != nil {
			i.ErrFn(path, err)
			return
		}

		contentChecksum := ChecksumPath(dst)
		i.ix.Add(pathChecksum, contentChecksum)
		return
	}

	// open the file
	f, err := os.Open(path)
	if err != nil {
		i.ErrFn(path, err)
		return
	}
	defer f.Close()

	// calculate content checksum
	contentChecksum, size, err := ChecksumReader(f)
	if err != nil {
		i.ErrFn(path, err)
		return
	}
	f.Close()

	if int64(size) != stat.Size() {
		i.ErrFn(path, fmt.Errorf("%w: %d vs %d", ErrChecksumSizeMismatch, size, stat.Size()))
		return
	}

	// append to the index
	i.ix.Add(pathChecksum, contentChecksum)
}

func (i *Indexer) processFiles(ctx context.Context) {
	defer i.wg.Done()

	for {
		select {
		case path := <-i.q:
			if path == "" {
				return
			}

			i.processFile(ctx, path)

		case <-ctx.Done():
			return
		}
	}
}

// Build walks the directory tree rooted at root and builds the index.
// It returns the index and an error if one occurs.
func (i *Indexer) Build(ctx context.Context) (*Index, error) {
	for range runtime.NumCPU() {
		i.wg.Add(1)
		go i.processFiles(ctx)
	}

	err := fs.WalkDir(os.DirFS(i.root), ".", func(path string, d os.DirEntry, err error) error {
		// report errors but continue walking
		if err != nil {
			i.ErrFn(path, err)
			return nil
		}

		// skip directories
		if d.IsDir() {
			return nil
		}

		// send for processing
		cleanPath, err := filepath.Abs(filepath.Join(i.root, path))
		i.q <- cleanPath

		return nil
	})

	close(i.q)
	i.wg.Wait()

	if err != nil {
		return nil, err
	}

	return i.ix, nil
}
