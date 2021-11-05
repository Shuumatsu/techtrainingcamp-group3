package sqlAPI

import "fmt"

type sqlApiError struct {
	FuncNotDefined error
}

var Error sqlApiError

func init() {
	Error.FuncNotDefined = fmt.Errorf("the function is not defined")
}
