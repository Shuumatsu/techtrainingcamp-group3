package sqlAPI

import (
	"techtrainingcamp-group3/db/dbmodels"
)

func FindOrCreateUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, SqlApiError.FuncNotDefined
}

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, SqlApiError.FuncNotDefined
}

func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	return nil, SqlApiError.FuncNotDefined
}

func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, SqlApiError.FuncNotDefined
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return SqlApiError.FuncNotDefined
}
