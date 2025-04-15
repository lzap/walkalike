package walkalike

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"io"
	"path/filepath"
	"time"

	"github.com/gnabgib/go-cksum"
)

func ChecksumPath(s string) uint32 {
	crc, length, err := cksum.Bytes([]byte(filepath.Clean(s)))

	if err != nil {
		panic(err)
	}

	if length != len(s) {
		panic("checksum length mismatch")
	}

	return crc
}

func ChecksumReader(r io.Reader) (uint32, int, error) {
	return cksum.Stream(bufio.NewReader(r))
}

func ChecksumCache(path string, size int64, modTime time.Time) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	h := crc64.New(crc64.MakeTable(crc64.ECMA))
	h.Write([]byte(absPath))
	h.Write([]byte(BuildID()))

	bSize := make([]byte, 8)
	binary.NativeEndian.PutUint64(bSize, uint64(size))
	h.Write(bSize)

	bMod := make([]byte, 8)
	binary.NativeEndian.PutUint64(bMod, uint64(modTime.UnixNano()))

	return fmt.Sprintf("%x-%s.index.bin", h.Sum64(), filepath.Base(path)), nil
}
