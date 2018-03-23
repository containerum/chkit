package update

import (
	"reflect"

	"io"

	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"
)

type Update struct {
	Binary    io.Reader `filename:"chkit"`
	Hash      io.Reader `filename:"sha256.sum"`
	Signature io.Reader `filename:"ecdsa.sig"`
}

func (u *Update) getFileMap() (ret map[string]int) {
	ret = make(map[string]int)
	for t, i := reflect.TypeOf(u), 0; i < t.NumField(); i++ {
		ret[t.Field(i).Tag.Get("filename")] = i
	}
	return
}

func (u *Update) Close() {
	for v, i := reflect.ValueOf(u), 0; i < v.NumField(); i++ {
		if cl, ok := v.Field(i).Interface().(io.Closer); ok {
			cl.Close()
		}
	}
}

const (
	ErrUnpack                = chkitErrors.Err("unable to unpack update file")
	ErrExpectedFilesNotFound = chkitErrors.Err("unable to find expected files in archive")
)

func unpack(rd io.Reader) (*Update, error) {
	var ret Update
	retVal := reflect.ValueOf(&ret)

	if err := unarchive(rd, &ret); err != nil {
		return nil, err
	}

	// check if we found all needed files in archive
	notFoundFiles := make([]string, 0)
	for i := 0; i < retVal.NumField(); i++ {
		if retVal.Field(i).Interface() == nil {
			notFoundFiles = append(notFoundFiles, retVal.Type().Field(i).Tag.Get("filename"))
		}
	}

	if len(notFoundFiles) > 0 {
		return nil, chkitErrors.Wrap(ErrExpectedFilesNotFound,
			fmt.Errorf("not found files:\n\t%s", strings.Join(notFoundFiles, "\n\t")))
	}

	return &ret, nil
}
