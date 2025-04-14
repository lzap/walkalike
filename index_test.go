package walkalike

import (
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

	t2 = []Token{
		{PathHash: 2, ContentHash: 2},
	}
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 0.01
}

func TestIntersectPaths(t *testing.T) {
	t.Parallel()

	a := t12
	b := t23
	expected := t2

	ix := &Index{Tokens: a}
	result := ix.IntersectPaths(b)

	if len(result) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(result))
	}

	for i, token := range result {
		if token != expected[i] {
			t.Errorf("Expected token %v, got %v", expected[i], token)
		}
	}
}

func TestIntersectContent(t *testing.T) {
	t.Parallel()

	a := t12
	b := t23
	expected := t2

	ix := &Index{Tokens: a}
	result := ix.IntersectContent(b)

	if len(result) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(result))
	}

	for i, token := range result {
		if token != expected[i] {
			t.Errorf("Expected token %v, got %v", expected[i], token)
		}
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

		result := aix.ContentSimilarityJaccard(bix)
		if !almostEqual(result, test.expected) {
			t.Errorf("Expected similarity %f, got %f", test.expected, result)
		}

		result = aix.PathSimilarityJaccard(bix)
		if !almostEqual(result, test.expected) {
			t.Errorf("Expected similarity %f, got %f", test.expected, result)
		}
	}
}
