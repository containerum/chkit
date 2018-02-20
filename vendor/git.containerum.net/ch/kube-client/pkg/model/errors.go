package model

import (
	"fmt"
)

type ResourceError struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"error"`
}

func (err *ResourceError) Error() string {
	return fmt.Sprintf("%s: %s ", err.Status, err.ErrorMsg)
}
