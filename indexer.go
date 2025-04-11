package walkalike

import (
	"context"
	"encoding/binary"
	"hash/crc64"
	"io/fs"
	"os"
	"sync"
)

type Indexer struct {
	// ErrFn is called when an error occurs during indexing. The first argument
	// is the path of the file that caused the error, and the second argument
	// is the error itself.
	ErrFn func(string, error)

	root fs.FS
	ix   *Index
	wg   *sync.WaitGroup
	q    chan item
}

type item struct {
	path string
	info os.DirEntry
}

// NewIndexer creates a new Indexer. The root argument is the root of the
// directory tree to be indexed.
func NewIndexer(root fs.FS) *Indexer {
	return &Indexer{
		ErrFn: func(path string, err error) {
			// Default error function does nothing
		},
		root: root,
		ix: &Index{
			Tokens: make([]Token, 0, 1024),
		},
		wg: &sync.WaitGroup{},
		q:  make(chan item, 1024),
	}
}

func (i *Indexer) processFiles(ctx context.Context) {
	defer i.wg.Done()

	emtyEntry := item{}

	// hash is shared since only one goroutine is writing to it
	table := crc64.MakeTable(crc64.ECMA)
	h := crc64.New(table)

	for {
		select {
		case entry := <-i.q:
			if entry == emtyEntry {
				return
			}

			if entry.info.IsDir() {
				continue
			}

			// prepare the hash
			h.Reset()

			// get the file info
			info, err := entry.info.Info()
			if err != nil {
				i.ErrFn(entry.path, err)
				continue
			}

			if info.IsDir() {
				continue
			}

			// write path to hash
			h.Write([]byte(entry.path))

			// write file size to hash
			if info.Size() > 0 {
				err := binary.Write(h, binary.NativeEndian, info.Size())
				if err != nil {
					i.ErrFn(entry.path, err)
					continue
				}
			} else {
				// zero-sized files are not indexed
				continue
			}

			// open the file
			f, err := i.root.Open(entry.path)
			if err != nil {
				i.ErrFn(entry.path, err)
				continue
			}

			// read the first block
			buf := make([]byte, 4096)
			n, err := f.Read(buf)
			if err != nil {
				i.ErrFn(entry.path, err)
				f.Close()
				continue
			}

			// write the first block to the hash
			if n > 0 {
				err := binary.Write(h, binary.NativeEndian, buf[:n])
				if err != nil {
					i.ErrFn(entry.path, err)
					f.Close()
					continue
				}
			}

			f.Close()

			// append to the index
			i.ix.Tokens = append(i.ix.Tokens, Token(h.Sum64()))

		case <-ctx.Done():
			return
		}
	}
}

// Build walks the directory tree rooted at root and builds the index.
// It returns the index and an error if one occurs. The index is built
// in lexicographic order of the file paths and keeps duplicate tokens.
func (i *Indexer) Build(ctx context.Context) (*Index, error) {
	// There must be exactly one goroutine that reads from the channel
	// because files must be processed in the order they are walked by
	// the deterministic WalkDir function.
	i.wg.Add(1)
	go i.processFiles(ctx)

	err := fs.WalkDir(i.root, ".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			i.ErrFn(path, err)
			// continue walking
			return nil
		}

		if d.IsDir() {
			return nil
		}

		item := item{
			path: path,
			info: d,
		}
		i.q <- item

		return nil
	})

	close(i.q)
	i.wg.Wait()

	if err != nil {
		return nil, err
	}

	return i.ix, nil
}
