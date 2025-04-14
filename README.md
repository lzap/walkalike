## walkalike - file tree similarity

Calculates how two or more filesystem trees are similar. Usage:

    go run github.com/lzap/walkalike/cmd dir1 dirN

Example output:

```
go run github.com/lzap/walkalike/cmd testdata/a testdata/b
0.44444444 testdata/a testdata/b
```

When multiple directories are provided, the first directory is compared against 2nd, 3rd and so on.

```
go run github.com/lzap/walkalike/cmd testdata/a testdata/b testdata/c
0.44444444 testdata/a testdata/b
0.96296296 testdata/a testdata/c
```

When a same directory is passed, the similarity is 1:

```
go run github.com/lzap/walkalike/cmd /usr/lib /usr/lib
1.00000000 /usr/lib /usr/lib
```

### Download binary

Use releases page on github to download binary for your OS and architecture.

### Files support

The tool support listing files inside OS images:

    sudo dnf -f install guestfs-tools

### Building

To build the project:

    go build ./cmd

### Output

Three columns are printed out:

* Similarity (number between 0 and 1)
* Directory A
* Directory B

### Implementation

The directory tree is walked and two CRC32 hashes (compatible with [cksum](https://github.com/coreutils/coreutils/blob/master/src/cksum.c)) are calculated for each file: absolute file path and file contents. Then the list of the CRC pairs is sorted by the file path hash ans stored as `index` which is used to perform similarity comparison.

Currently Jaccard similarity approach is used to calculate the final value. This may be configurable in the future as other algorithms are added.

### API

This repository is a Go library that can be used to achieve the same thing.

### TODO

* write guestfs wrapper for OS images
* inspect how this can be turned into RAG (Retrival Augmented Generation)
* validate the golang CRC32 is equal to `cksum` CRC 32 implementation
* 
