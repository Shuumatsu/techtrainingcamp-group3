package sqlAPI

import (
	"fmt"
	"techtrainingcamp-group3/db/dbmodels"
)

func FindUserByUID(uid dbmodels.UID) (dbmodels.User, error) {
	return dbmodels.User{}, fmt.Errorf("the function not defined")
}

func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	return nil, fmt.Errorf("the function is not defined")
}

func FindEnvelopeByEID(eid dbmodels.EID) (dbmodels.Envelope, error) {
	return dbmodels.Envelope{}, fmt.Errorf("the function not defined")
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return fmt.Errorf("the function not defined")
}
