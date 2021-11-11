package models

import "errors"

var ErrEnvelopeAlreadyOpen = errors.New("envelope already open")
var ErrEnvelopeOwner = errors.New("the envelope not belong to this user")
