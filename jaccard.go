package walkalike

import "sort"

// Intersect computes the intersection of two sets of tokens by cmp function.
// It returns a slice of tokens that are present in both sets.
// Both a and b arguments must be sorted and unique for the cmp function.
func Intersect(a, b []Token, cmp func(a Token, b Token) int) []Token {
	res := make([]Token, 0, max(len(a), len(b)))

	for i := range a {
		_, found := sort.Find(len(b), func(j int) int {
			return cmp(a[i], b[j])
		})

		if found {
			res = append(res, a[i])
		}
	}

	return res
}

func pathSimilarityJaccard(a, b *Index) float64 {
	as := SortByPaths(a.Tokens)
	bs := SortByPaths(b.Tokens)

	intersect := Intersect(as, bs, ComparePaths)
	union := len(as) + len(bs) - len(intersect)

	return float64(len(intersect)) / float64(union)
}

func contentSimilarityJaccard(a, b *Index) float64 {
	as := SortByContent(a.Tokens)
	bs := SortByContent(b.Tokens)

	intersect := Intersect(as, bs, CompareContent)
	union := len(as) + len(bs) - len(intersect)

	return float64(len(intersect)) / float64(union)
}

// JaccardSimilarity holds the Jaccard similarity between two indices.
// It contains the overall similarity, content similarity, and path similarity.
type JaccardSimilarity struct {
	// Similarity is the overall Jaccard similarity between two indices.
	Similarity        float64

	// ContentSimilarity is the Jaccard similarity between two indices by file content.
	ContentSimilarity float64

	// PathSimilarity is the Jaccard similarity between two indices by file paths.
	PathSimilarity    float64
}

// Similarity computes the Jaccard similarity between two indices.
//
// https://en.wikipedia.org/wiki/Jaccard_index
func SimilarityJaccard(a, b *Index) JaccardSimilarity {
	cs := contentSimilarityJaccard(a, b)
	ps := pathSimilarityJaccard(a, b)
	
	return JaccardSimilarity{
		Similarity:        (cs + ps) / 2,
		ContentSimilarity: cs,
		PathSimilarity:    ps,
	}
}
