package update

import (
	"archive/tar"
	"compress/gzip"
	"io"

	"reflect"

	"bytes"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

// unpack .tar.gz archive to temporary dir and save paths
func unarchive(rd io.Reader, update *Package) error {
	fmap := update.getFileMap()
	retVal := reflect.ValueOf(update)

	gzf, err := gzip.NewReader(io.LimitReader(rd, MaxFileSize))
	if err != nil {
		return chkitErrors.Wrap(ErrUnpack, err)
	}
	defer gzf.Close()

	tarf := tar.NewReader(gzf)

	for {
		header, nextErr := tarf.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return chkitErrors.Wrap(ErrUnpack, nextErr)
		}

		if header.Typeflag == tar.TypeReg {
			field, updateFile := fmap[header.Name]
			if updateFile {
				// this is our file, unpack it
				buf := bytes.NewBuffer(make([]byte, header.Size))

				if _, copyErr := io.Copy(buf, tarf); copyErr != nil {
					return chkitErrors.Wrap(ErrUnpack, copyErr)
				}

				retVal.Field(field).Set(reflect.ValueOf(io.Reader(buf)))
			}
		}
	}

	return nil
}
