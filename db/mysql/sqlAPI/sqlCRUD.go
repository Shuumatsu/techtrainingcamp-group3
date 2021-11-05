package sqlAPI

import (
	"fmt"
	"techtrainingcamp-group3/db/dbmodels"
)

type sqlAPIError struct {
	FuncNotDefined error
}

var SqlError sqlAPIError

func init() {
	SqlError.FuncNotDefined = fmt.Errorf("the function is not defined")
}

func FindOrCreateUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, SqlError.FuncNotDefined
}

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, SqlError.FuncNotDefined
}

func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	return nil, SqlError.FuncNotDefined
}

func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, SqlError.FuncNotDefined
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return SqlError.FuncNotDefined
}
