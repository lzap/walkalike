package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/lzap/walkalike"
)

func indexForDir(ctx context.Context, path string) (*walkalike.Index, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}
	root := os.DirFS(path)

	indexer := walkalike.NewIndexer(root)
	indexer.ErrFn = func(path string, err error) {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	index, err := indexer.Build(ctx)
	if err != nil {
		return nil, err
	}

	return index, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage: walkalike <path1> <pathN> ...")
		os.Exit(1)
	}

	indicies := make([]*walkalike.Index, 0, flag.NArg())
	for _, path := range flag.Args() {
		index, err := indexForDir(ctx, path)
		if err != nil {
			panic(err)
		}
		indicies = append(indicies, index)
	}

	for i := range indicies {
		println(indicies[i].String())
	}

	for i := 1; i < len(indicies); i++ {
		csim := indicies[0].ContentSimilarityJaccard(indicies[i])
		psim := indicies[0].PathSimilarityJaccard(indicies[i])
		fmt.Printf("%.13f %.13f %s %s\n", csim, psim, flag.Arg(0), flag.Arg(i))
	}
}
