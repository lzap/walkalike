package walkalike

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
)

type Index struct {
	Tokens []Token
}

type Token struct {
	PathHash    uint32
	ContentHash uint32
}

func init() {
	gob.Register(&Index{})
	gob.Register(&Token{})
}

// Add adds a new token to the index.
func (ix *Index) Add(pathHash, contentHash uint32) {
	ix.Tokens = append(ix.Tokens, Token{
		PathHash:    pathHash,
		ContentHash: contentHash,
	})
}

// Size returns the number of tokens in the index.
func (ix *Index) Size() int {
	return len(ix.Tokens)
}

// Encode encodes the index using gob and compresses it using gzip.
func (ix *Index) Encode(w io.Writer, filename string) error {
	nix := &Index{
		Tokens: SortByPaths(ix.Tokens),
	}

	gzw := gzip.NewWriter(w)
	defer gzw.Close()
	gzw.Name = filename
	gzw.Comment = "walkalike index v1"
	gzw.Extra = []byte{0x01} // format version

	return gob.NewEncoder(gzw).Encode(nix)
}

// Decode decodes the index using gob and decompresses it using gzip.
func (ix *Index) Decode(r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	return gob.NewDecoder(gzr).Decode(ix)
}

func (ix *Index) String() string {
	var s string
	for _, token := range ix.Tokens {
		s += token.String() + " "
	}
	return s
}

func (t Token) String() string {
	return fmt.Sprintf("%x:%x", t.PathHash, t.ContentHash)
}
