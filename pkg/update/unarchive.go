package update

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"

	"reflect"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

// unpack .tar.gz archive to temporary dir and save paths
func unarchive(rd io.Reader, tmpDir string, update *Update) error {
	fmap := update.getFileMap()
	retVal := reflect.ValueOf(update)

	gzf, err := gzip.NewReader(rd)
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
				fpath := path.Join(tmpDir, header.Name)
				file, createErr := os.Create(fpath)
				if createErr != nil {
					return chkitErrors.Wrap(ErrUnpack, createErr)
				}

				if _, copyErr := io.Copy(file, tarf); copyErr != nil {
					return chkitErrors.Wrap(ErrUnpack, copyErr)
				}

				retVal.Field(field).SetString(fpath)

				file.Close()
			}
		}
	}

	return nil
}
