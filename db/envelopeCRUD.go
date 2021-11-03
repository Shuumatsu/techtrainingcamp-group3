package db

import (
	"techtrainingcamp-group3/logger"
	"techtrainingcamp-group3/tools"
	"time"
	"strconv"
)

func GetEnvelope(envelope_id,user_id uint64) *Envelope{
	var envelope Envelope
	//get envelope according to envelope_id
	if err := DB.Table(Envelope{}.TableName()).First(&envelope,envelope_id).Error;err != nil{
		logger.Sugar.Warnw("GetEnvelope can not find envelope ","envelope_id",envelope_id)
		return nil
	}
	//check if user_id is right
	if envelope.Uid != user_id{
		logger.Sugar.Warnw("GetEnvelope envelope_id and user_id mismatch ","envelope_id",envelope_id,"user_id",user_id)
		return nil
	}

	return &envelope
}

func GetUser(user_id uint64) *User{
	var user User
	if err:= DB.Table(User{}.TableName()).First(&user,user_id).Error; err != nil{
		logger.Sugar.Warnw("open fail: cannot find user","user_id",user_id)
		return nil
	}
	return &user
}

func UpdateEnvelopeOpen(p *Envelope,u *User) error {
	//update envelope's open_stat and snatch_time
	tx := DB.Begin()
	if err:= tx.Table(Envelope{}.TableName()).Model(p).Update("open_stat", true).Error;err != nil{
		logger.Sugar.Errorw("UpdateEnvelopeOpen fail","envelope_id",p.Envelope_id)
		tx.Rollback()
		return err
	}
	//update user's amount added to envelope's value
	amountAfter := u.Amount + uint64(p.Value)
	if err:= tx.Table(User{}.TableName()).Where("uid", u.Uid).Update("amount", amountAfter).Error;err != nil{
		logger.Sugar.Errorw("UpdateEnvelopeOpen fail","envelope_id",p.Envelope_id)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UpdateUsersEnvelope(user User, cur_count int) (uint64, error) {
	envelope := tools.REPool.Snatch()
	// TODO: Rollback for red envelop pool
	eid := strconv.FormatUint(envelope.Eid, 10)
	tx := DB.Begin()
	if cur_count == 0 {
		if err := DB.Table(User{}.TableName()).Create(User{user.Uid, eid, 0}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		if err := DB.Table(User{}.TableName()).Where("uid", user.Uid).Update("envelope_list", user.EnvelopeList+","+eid).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := DB.Table(Envelope{}.TableName()).Create(Envelope{envelope.Eid, user.Uid, false, envelope.Money, time.Now().UTC()}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return envelope.Eid, nil
}