package walkalike

import (
	"encoding/csv"
	"io"
	"strconv"
)

type GuestfsLSCSV struct {
	r  io.Reader
	ix *Index
}

func NewGuestfsLSCSV(r io.Reader) *GuestfsLSCSV {
	return &GuestfsLSCSV{
		r: r,
		ix: &Index{
			Tokens: make([]Token, 0, 1024),
		},
	}
}

func (gr *GuestfsLSCSV) ReadAll() error {
	r := csv.NewReader(gr.r)
	r.FieldsPerRecord = -1

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// virt-ls output CSV format is as follows:
		//   type, mode, size, cksum, path, symlink_dest
		ftype := record[0]
		cksum := record[3]
		path := record[4]

		// skip directories or symlinks
		if ftype != "-" {
			continue
		}

		cksum64, err := strconv.ParseUint(cksum, 10, 32)
		if err != nil {
			return err
		}

		contentChecksum := uint32(cksum64)
		pathChecksum := ChecksumPath(path)

		gr.ix.Add(pathChecksum, contentChecksum)
	}

	return nil
}

func (gr *GuestfsLSCSV) Index() *Index {
	return gr.ix
}
