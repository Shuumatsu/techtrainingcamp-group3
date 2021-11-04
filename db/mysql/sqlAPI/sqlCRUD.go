package sqlAPI

import (
	"fmt"
	"techtrainingcamp-group3/db/dbmodels"
)

func FindUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	return nil, fmt.Errorf("the function not defined")
}

func FindEnvelopeByEID(eid dbmodels.EID) (*dbmodels.Envelope, error) {
	return nil, fmt.Errorf("the function not defined")
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return fmt.Errorf("the function not defined")
}
