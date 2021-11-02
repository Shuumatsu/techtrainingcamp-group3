package db

import (
	"techtrainingcamp-group3/logger"
	"time"
)

func GetEnvelope(envelope_id,user_id uint64) *Envelope{
	var envelope Envelope
	if err := DB.First(&envelope,envelope_id).Error;err != nil{
		logger.Sugar.Warnw("GetEnvelope can not find envelope ","envelope_id",envelope_id)
		return nil
	}
	if envelope.User_id != user_id{
		logger.Sugar.Warnw("GetEnvelope envelope_id and user_id mismatch ","envelope_id",envelope_id,"user_id",user_id)
		return nil
	}

	return &envelope
}

func UpdateEnvelopeOpen(p *Envelope) error {

	if err:= DB.Model(p).Updates(Envelope{Open_stat:true,Snatched_time:time.Now()}).Error;err != nil{
		logger.Sugar.Errorw("UpdateEnvelopeOpen fail","envelope_id",p.Envelope_id)
		return err
	}
	return nil
}