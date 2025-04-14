package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/lzap/walkalike"
)

var (
	Cache *walkalike.Cache
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

func indexForFile(ctx context.Context, path string) (*walkalike.Index, error) {
	if *verbose {
		fmt.Fprintf(os.Stderr, "Executing 'virt-ls' command")
	}
	cmd := exec.CommandContext(ctx, "virt-ls",
		"--csv",
		"--checksum=crc",
		"--long",
		"--recursive",
		"--add",
		path, "/",
	)
	cmd.Stderr = os.Stderr

	pr, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	guestfsBuilder := walkalike.NewGuestfsLSCSV(pr)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	err = guestfsBuilder.ReadAll()
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return guestfsBuilder.Index(), nil
}

func index(ctx context.Context, path string) (*walkalike.Index, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return indexForDir(ctx, path)
	} else {
		// check if we have a cached index
		ix, err := Cache.Get(path, info)
		if err != nil {
			return nil, err
		}
		if ix != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "Using cached index for %s (%d):\n", path, ix.Size())
			}
			return ix, nil
		}

		// if not, create a new index
		ix, err = indexForFile(ctx, path)
		if err != nil {
			return nil, err
		}

		// cache the index
		if *verbose {
			fmt.Fprintf(os.Stderr, "Built index for %s (%d):\n", path, ix.Size())
		}
		if err := Cache.Put(path, info, ix); err != nil {
			return nil, err
		}
		return ix, nil
	}
}

var verbose = flag.Bool("verbose", false, "enable verbose output")

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage: walkalike <path1> <pathN> ...")
		os.Exit(1)
	}

	Cache = walkalike.NewCache(filepath.Join(xdg.CacheHome, "walkalike"))
	err := Cache.Ensure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	indicies := make([]*walkalike.Index, 0, flag.NArg())
	for _, path := range flag.Args() {
		index, err := index(ctx, path)
		if err != nil {
			panic(err)
		}

		indicies = append(indicies, index)
	}

	for i := 1; i < len(indicies); i++ {
		sim := walkalike.SimilarityJaccard(indicies[0], indicies[i])
		if *verbose {
			fmt.Fprintf(os.Stderr, "%.13f %.13f %.13f %s %s\n",
				sim.Similarity,
				sim.ContentSimilarity,
				sim.PathSimilarity,
				flag.Arg(0), flag.Arg(i))
		}

		fmt.Printf("%.13f %s %s\n", sim.Similarity, flag.Arg(0), flag.Arg(i))
	}
}
