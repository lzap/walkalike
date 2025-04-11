package walkalike

import (
	"compress/gzip"
	"encoding/gob"
	"io"
	"slices"

	"github.com/s0rg/set"
)

type Index struct {
	Tokens []Token
}

type Token uint64

func init() {
	gob.Register(&Index{})
	gob.Register(Token(0))
}

// Uint64 returns the uint64 representation of the token.
func (t Token) Uint64() uint64 {
	return uint64(t)
}

// Comapact removes duplicate tokens from the index. It sorts the tokens
// before removing duplicates, so the order of the tokens is not preserved.
// This method is NOT called during indexing or encoding, the default encoder
// does store the index in the original order.
func (ix *Index) Compact() {
	if len(ix.Tokens) == 0 {
		return
	}

	slices.Sort(ix.Tokens)
	ix.Tokens = slices.Compact(ix.Tokens)
}

// Encode encodes the index using gob and compresses it using gzip.
func (ix *Index) Encode(w io.Writer) error {
	gzw := gzip.NewWriter(w)
	return gob.NewEncoder(gzw).Encode(ix)
}

// Decode decodes the index using gob and decompresses it using gzip.
func (ix *Index) Decode(r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	return gob.NewDecoder(gzr).Decode(ix)
}

// Similarity computes the Jaccard similarity between two indices.
//
// https://en.wikipedia.org/wiki/Jaccard_index
func (ix *Index) SimilarityJaccard(other *Index) float64 {
	a := make(set.Unordered[uint64])
	b := make(set.Unordered[uint64])

	for _, token := range ix.Tokens {
		a.Add(token.Uint64())
	}

	for _, token := range other.Tokens {
		b.Add(token.Uint64())
	}

	return float64(set.Intersect(a, b).Len()) / float64(set.Union(a, b).Len())
}
