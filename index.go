package walkalike

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"slices"
	"sort"
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

// ComparePaths compares two tokens by their path hash.
func ComparePaths(a, b Token) int {
	if a.PathHash > b.PathHash {
		return 1
	} else if a.PathHash < b.PathHash {
		return -1
	}
	return 0
}

// CompareContent compares two tokens by their content hash.
func CompareContent(a, b Token) int {
	if a.ContentHash > b.ContentHash {
		return 1
	} else if a.ContentHash < b.ContentHash {
		return -1
	}
	return 0
}

func SortByPaths(t []Token) []Token {
	res := make([]Token, len(t))
	copy(res, t)

	slices.SortFunc(res, func(i, j Token) int {
		return ComparePaths(i, j)
	})
	return res
}

func SortByContent(t []Token) []Token {
	res := make([]Token, len(t))
	copy(res, t)

	slices.SortFunc(res, func(i, j Token) int {
		return CompareContent(i, j)
	})
	return res
}

// Encode encodes the index using gob and compresses it using gzip.
func (ix *Index) Encode(w io.Writer) error {
	nix := &Index{
		Tokens: SortByPaths(ix.Tokens),
	}

	gzw := gzip.NewWriter(w)
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

// IntersectContent computes the intersection of two sets of tokens by file content.
// It returns a slice of tokens that are present in both sets.
// Index MUST be sorted before calling this function.
// The result is sorted in the same order as the original index.
func (ix *Index) IntersectContent(b []Token) []Token {
	a := ix.Tokens
	b = SortByContent(b)
	blen := len(b)
	res := make([]Token, 0, len(a))

	for i := range a {
		_, found := sort.Find(blen, func(j int) int {
			return CompareContent(a[i], b[j])
		})

		if found {
			res = append(res, a[i])
		}
	}

	return res
}

// IntersectPaths computes the intersection of two sets of tokens by file paths.
// It returns a slice of tokens that are present in both sets.
// Index MUST be sorted before calling this function.
// The result is sorted in the same order as the original index.
func (ix *Index) IntersectPaths(b []Token) []Token {
	a := ix.Tokens
	b = SortByPaths(b)
	blen := len(b)
	res := make([]Token, 0, len(a))

	for i := range a {
		_, found := sort.Find(blen, func(j int) int {
			return ComparePaths(a[i], b[j])
		})

		if found {
			res = append(res, a[i])
		}
	}

	return res
}

/*
func (ix *Index) UnionPaths(b []Token) []Token {
	a := ix.Tokens
	res := make([]Token, 0, len(a)+len(b))
	res = append(res, a...)
	res = append(res, b...)
	slices.SortFunc(res, func(i, j Token) int {
		return ComparePaths(i, j)
	})
	return slices.CompactFunc(res, func(i, j Token) bool {
		return ComparePaths(i, j) == 0
	})
}
*/

// Similarity computes the Jaccard similarity between two indices by file paths.
//
// https://en.wikipedia.org/wiki/Jaccard_index
func (ix *Index) PathSimilarityJaccard(other *Index) float64 {
	intersect := ix.IntersectPaths(other.Tokens)
	union := len(ix.Tokens) + len(other.Tokens) - len(intersect)

	return float64(len(intersect)) / float64(union)
}

// Similarity computes the Jaccard similarity between two indices by file content.
//
// https://en.wikipedia.org/wiki/Jaccard_index
func (ix *Index) ContentSimilarityJaccard(other *Index) float64 {
	intersect := ix.IntersectContent(other.Tokens)
	union := len(ix.Tokens) + len(other.Tokens) - len(intersect)

	return float64(len(intersect)) / float64(union)
}
