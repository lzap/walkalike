package walkalike

import (
	"fmt"
	"math"
	"testing"
)

var (
	t12 = []Token{
		{PathHash: 1, ContentHash: 1},
		{PathHash: 2, ContentHash: 2},
	}

	t23 = []Token{
		{PathHash: 2, ContentHash: 2},
		{PathHash: 3, ContentHash: 3},
	}

	t123 = []Token{
		{PathHash: 1, ContentHash: 1},
		{PathHash: 2, ContentHash: 2},
		{PathHash: 3, ContentHash: 3},
	}

	t1 = []Token{
		{PathHash: 1, ContentHash: 1},
	}

	t2 = []Token{
		{PathHash: 2, ContentHash: 2},
	}

	t3 = []Token{
		{PathHash: 3, ContentHash: 3},
	}
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 0.01
}

func TestIntersect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a, b, expected []Token
		cmp            func(a, b Token) int
	}{
		{t12, t23, t2, CompareContent},
		{t12, t123, t12, CompareContent},
		{t12, t2, t2, CompareContent},
		{t123, t12, t12, CompareContent},
		{t123, t23, t23, CompareContent},
		{t123, t1, t1, CompareContent},
		{t123, t2, t2, CompareContent},
		{t123, t3, t3, CompareContent},
		{t123, []Token{}, []Token{}, CompareContent},
		{t12, t23, t2, ComparePaths},
		{t12, t123, t12, ComparePaths},
		{t12, t2, t2, ComparePaths},
	}

	for ti, tst := range tests {
		t.Run(fmt.Sprintf("%d", ti), func(t *testing.T) {
			result := Intersect(tst.a, tst.b, CompareContent)

			if len(result) != len(tst.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tst.expected), len(result))
			}

			for i, token := range result {
				if token != tst.expected[i] {
					t.Errorf("Expected token %v, got %v", tst.expected[i], token)
				}
			}
		})
	}
}

func TestSimilarityJaccard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a, b     []Token
		expected float64
	}{
		{t12, t23, 0.33333333},
		{t12, t123, 0.66666666},
		{t12, t2, 0.5},
		{t23, t123, 0.66666666},
		{t23, t2, 0.5},
		{t123, t2, 0.33333333},
		{t2, t2, 1.0},
		{t12, t12, 1.0},
		{t123, t123, 1.0},
	}

	for _, test := range tests {
		aix := &Index{Tokens: test.a}
		bix := &Index{Tokens: test.b}
		result := SimilarityJaccard(aix, bix).Similarity
		
		if !almostEqual(result, test.expected) {
			t.Errorf("Expected similarity %f, got %f", test.expected, result)
		}
	}
}
