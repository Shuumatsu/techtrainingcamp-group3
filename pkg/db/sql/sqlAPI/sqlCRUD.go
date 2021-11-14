package sqlAPI

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"techtrainingcamp-group3/pkg/db/bloomfilter"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/kfk"
	"techtrainingcamp-group3/pkg/db/sql"
	"techtrainingcamp-group3/pkg/db/tokenBucket"
	"techtrainingcamp-group3/pkg/logger"
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
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doFindOrCreateUserByUID(defaultUser)
}
func doFindOrCreateUserByUID(defaultUser dbmodels.User) (*dbmodels.User, error) {
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).FirstOrCreate(
		&defaultUser).Error; err != nil {
		logger.Sugar.Debugw("FindOrCreateUserByUID", "error", err)
		return nil, err
	}
	bloomfilter.RedisAddUser(defaultUser.Uid)
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
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doFindUserByUID(uid)
}
func doFindUserByUID(uid dbmodels.UID) (*dbmodels.User, error) {
	user := dbmodels.User{Uid: uid}
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
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doFindEnvelopesByUID(uid)
}
func doFindEnvelopesByUID(uid dbmodels.UID) ([]dbmodels.Envelope, error) {
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
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return err
	}
	kfk.AddEnvelopeToUser(uid, envelope)
	return nil
	//return doAddEnvelopeToUserByUID(uid, envelope)
}
func doAddEnvelopeToUserByUID(uid dbmodels.UID, envelope dbmodels.Envelope) error {
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
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doFindEnvelopeByEID(eid)
}
func doFindEnvelopeByEID(eid dbmodels.EID) (*dbmodels.Envelope, error) {
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
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doOpenEnvelope(eid, uid)
}
func doOpenEnvelope(eid dbmodels.EID, uid dbmodels.UID) (*dbmodels.Envelope, error) {
	envelope, err := FindEnvelopeByEID(eid)
	if err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		return nil, err
	}
	if envelope.Opened == true {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		return nil, dbmodels.Error.EnvelopeAlreadyOpen
	}
	if envelope.Uid != uid {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		return nil, dbmodels.Error.ErrorEnvelopeOwner
	}
	kfk.OpenEnvelope(uid, *envelope)
	envelope.Opened = true
	return envelope, nil
}

// FindEnvelopeByUidEid
//	param:
//		eid > 0
//		uid > 0
//	return:
//		If cannot find eid in sql : Error.Notfound
//		else if the envelope's owner is not the param-uid: dbmodels.Error.ErrorEnvelopeOwner
func FindEnvelopeByUidEid(eid dbmodels.EID, uid dbmodels.UID) (*dbmodels.Envelope, error) {
	if eid == 0 || uid == 0 {
		return nil, Error.ErrorParam
	}
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doFindEnvelopeByUidEid(eid, uid)
}
func doFindEnvelopeByUidEid(eid dbmodels.EID, uid dbmodels.UID) (*dbmodels.Envelope, error) {
	var envelope dbmodels.Envelope
	// get envelope according to envelope_id
	if err := sql.DB.Table(dbmodels.Envelope{}.TableName()).Take(&envelope, eid).Error; err != nil {
		logger.Sugar.Warnw("GetEnvelope can not find envelope ", "envelope_id", eid)
		return nil, Error.NotFound
	}
	// check if user_id is right
	if envelope.Uid != uid {
		logger.Sugar.Warnw("GetEnvelope envelope_id and user_id mismatch ", "envelope_id", eid, "user_id", uid)
		return nil, dbmodels.Error.ErrorEnvelopeOwner
	}
	return &envelope, nil
}

// UpdateEnvelopeOpen
// param:
//		p *dbmodels.Envelope != nil
//		u *dbmodels.User != nil
// return:
//		if success : nil
//      if fail : database error
// update the envelope opened from false to true and add the envelope's value to user amount
func UpdateEnvelopeOpen(p *dbmodels.Envelope) (*dbmodels.User, error) {
	if p.EnvelopeId == 0 || p.Uid == 0 {
		return nil, Error.ErrorParam
	}
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return nil, err
	}
	return doUpdateEnvelopeOpen(p)
}

func doUpdateEnvelopeOpen(p *dbmodels.Envelope) (*dbmodels.User, error) {
	var user dbmodels.User
	var envelope dbmodels.Envelope
	user.Uid = p.Uid
	envelope.EnvelopeId = p.EnvelopeId

	// find envelope status
	if err := sql.DB.Table(dbmodels.Envelope{}.TableName()).Take(&envelope).Error; err != nil {
		logger.Sugar.Errorw("Find Envelope By EID", "error", err)
		return nil, err
	}

	// check open status for data consistency
	if envelope.Opened == true {
		return nil, dbmodels.Error.EnvelopeAlreadyOpen
	}

	// find user amount to get amount
	if err := sql.DB.Table(dbmodels.User{}.TableName()).Take(&user).Error; err != nil {
		logger.Sugar.Errorw("FindUserByUID", "error", err)
		return nil, err
	}

	kfk.OpenEnvelope(p.Uid, *p)

	user.Amount += p.Value

	return &user, nil
}
