package sqlAPI

import (
	"techtrainingcamp-group3/db/dbmodels"
)

func FindOrCreateUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, Error.FuncNotDefined
}

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, Error.FuncNotDefined
}

func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	return nil, Error.FuncNotDefined
}

func AddEnvelopeToUserByUID(uid dbmodels.UID, envelope dbmodels.Envelope) error {
	return Error.FuncNotDefined
}

func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, Error.FuncNotDefined
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return Error.FuncNotDefined
}
