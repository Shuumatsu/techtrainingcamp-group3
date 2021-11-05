package sqlAPI

import "fmt"

type sqlApiError struct {
	FuncNotDefined error
}

var SqlApiError sqlApiError

func init() {
	SqlApiError.FuncNotDefined = fmt.Errorf("the function is not defined")
}
