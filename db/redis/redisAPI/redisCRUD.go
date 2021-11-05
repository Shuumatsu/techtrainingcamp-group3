package redisAPI

import (
	"fmt"
	"techtrainingcamp-group3/db/dbmodels"
)

type redisAPIError struct {
	FuncNotDefined   error
	UserNotExist     error
	EnvelopeNotExist error
}

var RdeisError redisAPIError

func init() {
	RdeisError.FuncNotDefined = fmt.Errorf("the function is not defined")
}

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, RdeisError.FuncNotDefined
}

func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, RdeisError.FuncNotDefined
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return RdeisError.FuncNotDefined
}
