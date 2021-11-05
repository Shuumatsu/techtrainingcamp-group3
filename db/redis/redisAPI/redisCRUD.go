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

var RedisError redisAPIError

func init() {
	RedisError.FuncNotDefined = fmt.Errorf("the function is not defined")
}

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, RedisError.FuncNotDefined
}

func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, RedisError.FuncNotDefined
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return RedisError.FuncNotDefined
}
