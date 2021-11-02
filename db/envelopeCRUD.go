package db

import (
	"techtrainingcamp-group3/logger"
	"time"
)

func GetEnvelope(envelope_id,user_id uint64) *Envelope{
	var envelope Envelope
	//get envelope according to envelope_id
	if err := DB.Table(Envelope{}.TableName()).First(&envelope,envelope_id).Error;err != nil{
		logger.Sugar.Warnw("GetEnvelope can not find envelope ","envelope_id",envelope_id)
		return nil
	}
	//check if user_id is right
	if envelope.User_id != user_id{
		logger.Sugar.Warnw("GetEnvelope envelope_id and user_id mismatch ","envelope_id",envelope_id,"user_id",user_id)
		return nil
	}

	return &envelope
}

func GetUser(user_id uint64) *OpenUser{
	var openUser OpenUser
	if err:= DB.Table(OpenUser{}.TableName()).First(&openUser,user_id).Error; err != nil{
		logger.Sugar.Warnw("open fail: cannot find user","user_id",user_id)
		return nil
	}
	return &openUser
}

func UpdateEnvelopeOpen(p *Envelope,u *OpenUser) error {
	//update envelope's open_stat and snatch_time
	tx := DB.Begin()
	if err:= tx.Table(Envelope{}.TableName()).Model(p).Updates(Envelope{Open_stat:true,Snatched_time:time.Now()}).Error;err != nil{
		logger.Sugar.Errorw("UpdateEnvelopeOpen fail","envelope_id",p.Envelope_id)
		tx.Rollback()
		return err
	}
	//update user's amount added to envelope's value
	amountAfter := u.Amount + uint64(p.Value)
	if err:= tx.Table(OpenUser{}.TableName()).Model(u).Update("amount", amountAfter).Error;err != nil{
		logger.Sugar.Errorw("UpdateEnvelopeOpen fail","envelope_id",p.Envelope_id)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}