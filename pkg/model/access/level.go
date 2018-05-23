package access

import (
	"encoding"

	"github.com/containerum/chkit/pkg/chkitErrors"
	"github.com/containerum/chkit/pkg/util/strset"
)

type AccessLevel string

const (
	Owner      AccessLevel = "owner"
	Write      AccessLevel = "write"
	ReadDelete AccessLevel = "read-delete"
	Read       AccessLevel = "read"
	None       AccessLevel = ""

	ErrInvalidAccessLevel chkitErrors.Err = "invalid access level"
)

var (
	_ encoding.TextMarshaler   = None
	_ encoding.TextUnmarshaler = new(AccessLevel)
)

var levelsList = []AccessLevel{
	Owner,
	Write,
	ReadDelete,
	Read,
	None,
}

func Levels() []AccessLevel {
	return append(make([]AccessLevel, 0, len(levelsList)), levelsList...)
}

var lvls = strset.NewSet([]string{
	Owner.String(),
	Write.String(),
	ReadDelete.String(),
	Read.String(),
})

func LevelFromString(str string) (AccessLevel, error) {
	if str == "" {
		return None, nil
	}
	if lvls.Have(str) {
		return AccessLevel(str), nil
	}
	return None, ErrInvalidAccessLevel.Comment(str)
}

func (lvl AccessLevel) String() string {
	if None == lvl {
		return "none"
	}
	return string(lvl)
}

func (lvl AccessLevel) MarshalText() ([]byte, error) {
	return []byte(lvl.String()), nil
}

func (lvl *AccessLevel) UnmarshalText(p []byte) error {
	l, err := LevelFromString(string(p))
	if err != nil {
		return err
	}
	*lvl = l
	return nil
}
