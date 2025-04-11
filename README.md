## walkalike - file tree similarity

Calculates how two or more filesystem trees are similar. Usage:

    go run github.com/lzap/walkalike@latest dir1 dir2

Example output:

```
go run github.com/lzap/walkalike@latest testdata/a testdata/b testdata/c
0.44444444 testdata/a testdata/b
0.96296296 testdata/a testdata/c
```

### Download binary

Use releases page on github to download binary for your OS and architecture.

### Output

Three columns are printed out:

* Similarity (number between 0 and 1)
* Directory A
* Directory B

When multiple directories are provided, the first directory is compared against 2nd, 3rd and so on.

### Implementation

The directory tree is walked in lexicographic order and CRC64 is calculated for each file path, size and the initial 4096 bytes of data. This list of CRC64 values, also called `index`, is used to perform similarity comparison.

Currently Jaccard similarity approach is used to calculate the final value. This may be configurable in the future as other algorithms are added.

### API

This repository is a Go library that can be used to achieve the same thing.

### TODO

* write guestfs wrapper for OS images
* inspect how this can be turned into RAG (Retrival Augmented Generation)
