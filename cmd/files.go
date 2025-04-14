package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lzap/walkalike"
)

var indexForFile = func(ctx context.Context, path string) (*walkalike.Index, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is not a file", path)
	}
	_ = ctx

	return &walkalike.Index{}, nil
}
