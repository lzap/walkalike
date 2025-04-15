## walkalike - file tree similarity

Calculates how two or more filesystem trees (or OS images) are similar by applying techniques from Machine Learning, Deep Learning, Big Data and Clustering. In practice, the tool prints out Jaccard similarity coefficient for two directories. Usage:

    go run github.com/lzap/walkalike/cmd dir1 dirN

Processing is optimized for SSDs and index creation is parallelized. Example command:

    go run github.com/lzap/walkalike/cmd testdata/a testdata/b

Should print:

    0.4444444444444 testdata/a testdata/b

Similarity coefficient is between 0.0, when two trees are not similar at all, or 1.0, when two trees are exactly the same. When multiple directories are provided, the first directory is compared against 2nd, 3rd and so on.

    go run github.com/lzap/walkalike/cmd testdata/a testdata/b testdata/c

Should print:

    0.4444444444444 testdata/a testdata/b
    0.9814814814815 testdata/a testdata/c

When the same directory is passed twice:

    go run github.com/lzap/walkalike/cmd /usr/lib /usr/lib

Similarity must be 1.0:

    1.0000000000000 /usr/lib /usr/lib

### OS images support

The tool support listing files inside OS images, it requires the tool `virt-ls` to be installed:

    sudo dnf -f install guestfs-tools

The utility from guestfs-tools is executed as a subprocess and automatically detect OS image, finds root partition and walks the directory tree passing results and checksums to the parent process which calculates the index.

    go run github.com/lzap/walkalike/cmd a-fedora-40-minimal-raw-x86_64.raw b-fedora-40-minimal-raw-x86_64.raw

Building the indices will take a little bit more time as walking the tree is not parallelized and some images might be compressed but the result is exactly the same:

    0.9995304751005 a-fedora-40-minimal-raw-x86_64.raw b-fedora-40-minimal-raw-x86_64.raw

### Output format

Three columns are printed out:

* Similarity (number between 0 and 1)
* Directory/OS image A
* Directory/OS image B

### Index cache

Indices for OS images are kept in cache (`$HOME/.cache/walkalike`) so further processing will be fast. Directories are never kept in cache however.

Index size is relatively small, for an OS image of about 35k files about 500kB index is created. The rough estimation is 16 bytes per file entry, files are compressed with gzip which brings the size down by about 20%.

### Implementation

When a directory is passed, the tree is walked and two CRC32 checksums (compatible with [cksum](https://github.com/coreutils/coreutils/blob/master/src/cksum.c)) are calculated for each file: absolute file path and file contents. Then the list of the CRC pairs is stored in XDG_CACHE directory as a small `index` which is used to perform similarity comparison.

When a file is passed, `virt-ls` tool is called to detect OS image type, partitions and list all files including the CRC32 checksum. The same checksum is calculated from the file paths as well.

Only regular files, hardlinks and symlinks are subject of processing, directories are skipped. Meaning, an empty directory will not affect the result score.

Currently [Jaccard similarity](https://en.wikipedia.org/wiki/Jaccard_index) approach is used to calculate the final value. Calculation is done separately for file paths and content checksums with the weight of 0.5. This may be configurable in the future as other algorithms are added.

### Download binary

Use releases page on github to download binary for your OS and architecture.

### Building

To build the project:

    go build ./cmd

### API

This repository is also a Go library.

### Authors

(c) 2025 Lukáš Zapletal

### License

Apache 2.0
