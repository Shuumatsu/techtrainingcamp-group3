package sqlAPI

import (
	"fmt"
	"gorm.io/gorm"
)

type sqlApiError struct {
	FuncNotDefined      error
	NotFound            error
	ErrorParam          error
	EnvelopeAlreadyOpen error
	ErrorEnvelopeOwner  error
}

var Error sqlApiError

func init() {
	Error.FuncNotDefined = fmt.Errorf("the function is not defined")
	Error.NotFound = gorm.ErrRecordNotFound
	Error.ErrorParam = fmt.Errorf("the param is error")
	Error.EnvelopeAlreadyOpen = fmt.Errorf("the envelope already open")
	Error.ErrorEnvelopeOwner = fmt.Errorf("the envelope not belong to this user")
}
