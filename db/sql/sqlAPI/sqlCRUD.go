package sqlAPI

import (
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql"
	"techtrainingcamp-group3/logger"
)

func FindOrCreateUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	var user dbmodels.User
	if err := sql.DB.
		Table(dbmodels.User{}.TableName()).
		FirstOrCreate(&user, uid).Error; err != nil {
		logger.Sugar.Debugw("FindOrCreateUserByUID", "error", err)
		return nil, err
	}
	return &user, nil
}

func FindUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	var user dbmodels.User
	if err := sql.DB.
		Table(dbmodels.User{}.TableName()).
		First(&user, uid).Error; err != nil {
		logger.Sugar.Debugw("FindUserByUID", "error", err)
		return nil, err
	}
	return &user, nil
}

func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	return nil, Error.FuncNotDefined
}

func AddEnvelopeToUserByUID(uid dbmodels.UID, envelope dbmodels.Envelope) error {
	return Error.FuncNotDefined
}

func FindEnvelopeByEID(eid dbmodels.EID) (*dbmodels.Envelope, error) {
	var envelope dbmodels.Envelope
	if err := sql.DB.
		Table(dbmodels.Envelope{}.TableName()).
		First(&envelope, eid).Error; err != nil {
		logger.Sugar.Debugw("FindEnvelopeByEID", "error", err)
		return nil, err
	}
	return &envelope, nil
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return Error.FuncNotDefined
}
