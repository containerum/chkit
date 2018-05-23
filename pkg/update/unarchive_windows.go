package update

import (
	"archive/zip"
	"io"

	"bytes"
	"io/ioutil"

	"reflect"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

// unpack .zip archive and save paths
func unarchive(rd io.Reader, update *Package) error {
	fmap := update.getFileMap()
	retVal := reflect.ValueOf(update).Elem()

	// zip file requires random access so before start we read all "rd" contents to buffer
	content, err := ioutil.ReadAll(io.LimitReader(rd, MaxFileSize))
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
				rc, openErr := archFile.Open()
				if openErr != nil {
					return chkitErrors.Wrap(ErrUnpack, openErr)
				}

				buf := bytes.NewBuffer(make([]byte, archFile.UncompressedSize64))

				if _, copyErr := io.Copy(buf, rc); copyErr != nil {
					return chkitErrors.Wrap(ErrUnpack, copyErr)
				}

				retVal.Field(field).Set(reflect.ValueOf(io.Reader(buf)))

				rc.Close()
			}
		}
	}

	return nil
}
