package dbmodels

import "fmt"

type dbError struct{
	EnvelopeAlreadyOpen error
	ErrorEnvelopeOwner  error
}
var Error dbError

func init()  {
	Error.EnvelopeAlreadyOpen = fmt.Errorf("the envelope already open")
	Error.ErrorEnvelopeOwner = fmt.Errorf("the envelope not belong to this user")
}