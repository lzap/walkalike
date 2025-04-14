package walkalike

import "slices"

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

// SortByPaths sorts the tokens by their path hash.
// It returns a new slice of tokens with duplicates removed.
func SortByPaths(t []Token) []Token {
	res := make([]Token, len(t))
	copy(res, t)

	slices.SortFunc(res, func(i, j Token) int {
		return ComparePaths(i, j)
	})

	return slices.CompactFunc(res, func(i, j Token) bool {
		return i.PathHash == j.PathHash
	})
}

// SortByContent sorts the tokens by their content hash.
// It returns a new slice of tokens with duplicates removed.
func SortByContent(t []Token) []Token {
	res := make([]Token, len(t))
	copy(res, t)

	slices.SortFunc(res, func(i, j Token) int {
		return CompareContent(i, j)
	})

	return slices.CompactFunc(res, func(i, j Token) bool {
		return i.ContentHash == j.ContentHash
	})
}
