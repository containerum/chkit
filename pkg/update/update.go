package update

import (
	"os"
	"path"
	"reflect"

	"io"

	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

type Update struct {
	BinaryPath    string `filename:"chkit"`
	HashPath      string `filename:"sha256.sum"`
	SignaturePath string `filename:"ecdsa.sig"`
}

func (u *Update) getFileMap() (ret map[string]int) {
	ret = make(map[string]int)
	for t, i := reflect.TypeOf(u), 0; i < t.NumField(); i++ {
		ret[t.Field(i).Tag.Get("filename")] = i
	}
	return
}

const (
	ErrUnpack                = chkitErrors.Err("unable to unpack update file")
	ErrExpectedFilesNotFound = chkitErrors.Err("unable to find expected files in archive")
)

func unpack(rd io.Reader) (*Update, error) {
	tmpDir := path.Join(os.TempDir(), "containerum")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil && !os.IsExist(err) {
		return nil, chkitErrors.Wrap(ErrUnpack, err)
	}

	var ret Update
	retVal := reflect.ValueOf(&ret)

	if err := unarchive(rd, tmpDir, &ret); err != nil {
		return nil, err
	}

	// check if we found all needed files in archive
	notFoundFiles := make([]string, 0)
	for i := 0; i < retVal.NumField(); i++ {
		if retVal.Field(i).String() == "" {
			notFoundFiles = append(notFoundFiles, retVal.Type().Field(i).Tag.Get("filename"))
		}
	}

	if len(notFoundFiles) > 0 {
		return nil, chkitErrors.Wrap(ErrExpectedFilesNotFound,
			fmt.Errorf("not found files:\n\t%s", strings.Join(notFoundFiles, "\n\t")))
	}

	return &ret, nil
}
