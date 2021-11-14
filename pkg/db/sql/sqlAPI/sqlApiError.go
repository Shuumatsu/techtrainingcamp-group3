package sqlAPI

import (
	"fmt"
	"gorm.io/gorm"
)

type sqlApiError struct {
	FuncNotDefined      error
	NotFound            error
	ErrorParam          error
}

var Error sqlApiError

func init() {
	Error.FuncNotDefined = fmt.Errorf("the function is not defined")
	Error.NotFound = gorm.ErrRecordNotFound
	Error.ErrorParam = fmt.Errorf("the param is error")
}
