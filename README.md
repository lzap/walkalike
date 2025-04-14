## walkalike - file tree similarity

Calculates how two or more filesystem trees are similar. Usage:

    go run github.com/lzap/walkalike/cmd dir1 dirN

Example output:

```
go run github.com/lzap/walkalike/cmd testdata/a testdata/b
0.4444444444444 testdata/a testdata/b
```

When multiple directories are provided, the first directory is compared against 2nd, 3rd and so on.

```
go run github.com/lzap/walkalike/cmd testdata/a testdata/b testdata/c
0.4444444444444 testdata/a testdata/b
0.9814814814815 testdata/a testdata/c
```

When a same directory is passed, the similarity is 1:

```
go run github.com/lzap/walkalike/cmd /usr/lib /usr/lib
1.0000000000000 /usr/lib /usr/lib
```

### OS images support

The tool support listing files inside OS images, it requires the tool `virt-ls` to be installed:

    sudo dnf -f install guestfs-tools

```
0.9995304751005 a-fedora-40-minimal-raw-x86_64.raw b-fedora-40-minimal-raw-x86_64.raw
```

### Building

To build the project:

    go build ./cmd

### Download binary

Use releases page on github to download binary for your OS and architecture.

### Output

Three columns are printed out:

* Similarity (number between 0 and 1)
* Directory A
* Directory B

### Implementation

When a directory is passed, the tree is walked and two CRC32 checksums (compatible with [cksum](https://github.com/coreutils/coreutils/blob/master/src/cksum.c)) are calculated for each file: absolute file path and file contents. Then the list of the CRC pairs is stored in XDG_CACHE directory as a small `index` which is used to perform similarity comparison.

When a file is passed, `virt-ls` tool is called to detect OS image type, partitions and list all files including the CRC32 checksum. The same checksum is calculated from the file paths as well.

Directories are not stored in the index. Symlinks are not followed.

Currently Jaccard similarity approach is used to calculate the final value. Calculation is done separately for file paths and content checksums with the weight of 0.5. This may be configurable in the future as other algorithms are added.

### LLVM RAG

[TODO](https://en.wikipedia.org/wiki/Retrieval-augmented_generation)

### API

This repository is also a Go library.
