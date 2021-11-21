package sqlAPI

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"techtrainingcamp-group3/pkg/db/bloomfilter"
	"techtrainingcamp-group3/pkg/db/dbmodels"
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
	return doAddEnvelopeToUserByUID(uid, envelope)
}
func doAddEnvelopeToUserByUID(uid dbmodels.UID, envelope dbmodels.Envelope) error {
	//start transaction
	tx := sql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//Create envelope in envelope table
	if err := tx.Table(
		dbmodels.Envelope{}.TableName()).Create(
		&envelope).Error; err != nil {
		logger.Sugar.Debugw("AddEnvelopeToUserByUID", "error", err)
		tx.Rollback()
		return err
	}

	//add envelope to user's envelope list
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
func UpdateEnvelopeOpen(p *dbmodels.Envelope) error {
	if p.EnvelopeId == 0 || p.Uid == 0 {
		return Error.ErrorParam
	}
	err := tokenBucket.Limiter.Wait(context.TODO())
	if err != nil {
		return err
	}
	return doUpdateEnvelopeOpen(p)
}
func doUpdateEnvelopeOpen(p *dbmodels.Envelope) error {
	var user dbmodels.User
	var envelope dbmodels.Envelope
	user.Uid = p.Uid
	envelope.EnvelopeId = p.EnvelopeId
	//start transaction
	tx := sql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// find envelope status
	if err := sql.DB.Table(dbmodels.Envelope{}.TableName()).Take(&envelope).Error; err != nil {
		tx.Rollback()
		logger.Sugar.Errorw("Find Envelope By EID", "error", err)
		return err
	}
	// check open status for data consistency
	if envelope.Opened == true {
		tx.Rollback()
		return dbmodels.Error.EnvelopeAlreadyOpen
	}

	logger.Sugar.Debugw("Consumer: OpenEnvelopeByEID", "uid", p.Uid, "envelope", envelope)

	//Update envelope to be opened
	if err := tx.Model(
		&envelope).Update("opened", true).Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return err
	}

	//Update user amount
	if err := tx.Table(
		dbmodels.User{}.TableName()).Where("uid", p.Uid).Update(
		"amount", gorm.Expr(
			"amount + ?", envelope.Value)).Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return err
	}
	err := tx.Commit().Error
	if err != nil{
		logger.Sugar.Errorw("UpdateEnvelopeOpen sql error","eid",p.EnvelopeId)
		return err
	}

	return nil
}
