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
//		defaultUser(uid > 0): if not found, will create it
// 	return:
//		check by Uid
// 		if the user is exist, return the user;
//		else if the user is notFound, create the defaultuser and return it;
//		else return error(other databaseError)
func FindOrCreateUserByUID(defaultUser dbmodels.User) (*dbmodels.User, error) {
	if defaultUser.Uid == 0 {
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
	if uid == 0 {
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
// 		return all envelopes belong to the user
//		or return the error(other database error)
func FindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
	if uid == 0 {
		return nil, Error.ErrorParam
	}
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
//		envelope(envelopeId > 0): new envelope
// 	return:
//		error(other database error)
// create the envelope in envelope table and append it to the user's envelope_list
func AddEnvelopeToUserByUID(uid dbmodels.UID, envelope dbmodels.Envelope) error {
	if uid == 0 || envelope.EnvelopeId == 0 {
		return Error.ErrorParam
	}
	tx := sql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Table(
		dbmodels.Envelope{}.TableName()).Create(
		&envelope).Error; err != nil {
		logger.Sugar.Debugw("AddEnvelopeToUserByUID", "error", err)
		tx.Rollback()
		return err
	}
	if err := tx.Model(
		&dbmodels.User{Uid: uid}).Update(
		"envelope_list", gorm.Expr(
			fmt.Sprintf(`CONCAT(envelope_list,",%s")`, envelope.EnvelopeId.String()))).
		Error; err != nil {
		logger.Sugar.Debugw("AddEnvelopeToUserByUID", "error", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// FindEnvelopeByEID
//	param:
//		eid > 0
//	return:
// 		if the envelope is exist, return the envelope
//		else return error(Error.NotFound or other databaseError)
func FindEnvelopeByEID(eid dbmodels.EID) (*dbmodels.Envelope, error) {
	if eid == 0 {
		return nil, Error.ErrorParam
	}
	envelope := dbmodels.Envelope{EnvelopeId: eid}
	if err := sql.DB.Table(
		dbmodels.Envelope{}.TableName()).Take(
		&envelope).Error; err != nil {
		logger.Sugar.Debugw("FindEnvelopeByEID", "error", err)
		return nil, err
	}
	return &envelope, nil
}

// OpenEnvelope
//	param:
//		eid > 0
//		uid > 0
//	return:
//		if eid == 0 or uid == 0: return Error.ErrorParam
//		else if the envelope is already open: return Error.EnvelopeAlreadyOpen
//		else if the envelope's owner is not the param-uid: return Error.ErrorEnvelopeOwner
//		else if sql error: return other database error
//		else return the envelope
// update the envelope open and add the envelope's value to user amount
func OpenEnvelope(eid dbmodels.EID, uid dbmodels.UID) (*dbmodels.Envelope, error) {
	if eid == 0 || uid == 0 {
		return nil, Error.ErrorParam
	}
	tx := sql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// check the envelope
	envelope, err := FindEnvelopeByEID(eid)
	if err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return nil, err
	}
	if envelope.Opened == true {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return nil, Error.EnvelopeAlreadyOpen
	}
	if envelope.Uid != uid {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return nil, Error.ErrorEnvelopeOwner
	}
	// set envelope open
	if err := tx.Model(
		&envelope).Update("opened", true).Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return nil, err
	}
	logger.Sugar.Debugw("OpenEnvelopeByEID", "envelope", envelope)
	// add user amount
	user := dbmodels.User{Uid: uid}
	if err := tx.Model(
		&user).Update(
		"amount", gorm.Expr(
			"amount + ?", envelope.Value)).Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return nil, err
	}
	return envelope, nil
}
