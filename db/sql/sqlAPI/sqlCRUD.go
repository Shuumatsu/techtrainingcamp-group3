package sqlAPI

import (
	"fmt"
	"gorm.io/gorm"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql"
	"techtrainingcamp-group3/logger"
)

// FindOrCreateUserByUID
// 	param:
// 		uid > 0
//		defaultUser: if not found, will create it
// 	return:
// 		if the user is exist, return the user;
//		else if the user is notFound, create the defaultuser and return it;
//		else return error(other databaseError)
func FindOrCreateUserByUID(uid dbmodels.UID, defaultUser dbmodels.User) (*dbmodels.User, error) {
	if uid <= 0 {
		return nil, Error.ErrorParam
	}
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).FirstOrCreate(
		&defaultUser).Error; err != nil {
		logger.Sugar.Debugw("FindOrCreateUserByUID", "error", err)
		return nil, err
	}
	return &defaultUser, nil
}

// FindUserByUID
// 	param:
// 		uid > 0
// 	return:
// 		if the user is exist, return the user
//		else return error(Error.NotFound or other databaseError)
func FindUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	if uid <= 0 {
		return nil, Error.ErrorParam
	}
	var user dbmodels.User
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).Take(
		&user).Error; err != nil {
		logger.Sugar.Debugw("FindUserByUID", "error", err)
		return nil, err
	}
	return &user, nil
}

// FindEnvelopesByUID
// 	param:
// 		uid > 0
// 	return:
// 		all envelopes belong to the user
func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	var Envelopes []dbmodels.Envelope
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Where(
		"uid = ?", uid).
		Find(&Envelopes).Error; err != nil {
		return nil, err
	}
	return Envelopes, nil
}

// AddEnvelopeToUserByUID
// 	param:
// 		uid > 0
//		envelope: new envelope
// 	return:
//		error(other database error)
// create the envelope in envelope table and append it to the user's envelope_list
func AddEnvelopeToUserByUID(uid dbmodels.UID, envelope dbmodels.Envelope) error {
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Create(
		&envelope).Error; err != nil {
		logger.Sugar.Debugw("AddEnvelopeToUserByUID", "error", err)
		return err
	}
	if err := sql.DB.Model(
		&dbmodels.User{Uid: uid}).Update(
		"envelope_list", gorm.Expr(
			fmt.Sprintf(`CONCAT(envelope_list,",%s")`, envelope.EnvelopeId.String()))).
		Error; err != nil {
		logger.Sugar.Debugw("AddEnvelopeToUserByUID", "error", err)
		return err
	}
	return nil
}

func FindEnvelopeByEID(eid dbmodels.EID) (*dbmodels.Envelope, error) {
	var envelope dbmodels.Envelope
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Take(
		&envelope, eid).Error; err != nil {
		logger.Sugar.Debugw("FindEnvelopeByEID", "error", err)
		return nil, err
	}
	return &envelope, nil
}

func OpenEnvelopeByEID(eid dbmodels.EID) error {
	return Error.FuncNotDefined
}
