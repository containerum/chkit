package update

import (
	"archive/zip"
	"io"
	"os"
	"path"

	"bytes"
	"io/ioutil"

	"reflect"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

// unpack .zip archive and save paths
func unarchive(rd io.Reader, tmpDir string, update *Update) error {
	fmap := update.getFileMap()
	retVal := reflect.ValueOf(update)

	// zip file requires random access so before start we read all "rd" contents to buffer
	content, err := ioutil.ReadAll(rd)
	if err != nil {
		return chkitErrors.Wrap(ErrUnpack, err)
	}
	zipf, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return chkitErrors.Wrap(ErrUnpack, err)
	}

	for _, archFile := range zipf.File {
		if archFile.FileInfo().Mode().IsRegular() {
			field, updateFile := fmap[archFile.Name]
			if updateFile {
				// this is our file, unpack it
				fpath := path.Join(tmpDir, archFile.Name)
				file, createErr := os.Create(fpath)
				if createErr != nil {
					return chkitErrors.Wrap(ErrUnpack, createErr)
				}

				rc, openErr := archFile.Open()
				if openErr != nil {
					return chkitErrors.Wrap(ErrUnpack, openErr)
				}

				if _, copyErr := io.Copy(file, rc); copyErr != nil {
					return chkitErrors.Wrap(ErrUnpack, copyErr)
				}

				retVal.Field(field).SetString(fpath)

				rc.Close()
				file.Close()
			}
		}
	}

	return nil
}
