package kfk

import (
	"github.com/Shopify/sarama"
	"strconv"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/db/rds/redisAPI"
	"techtrainingcamp-group3/pkg/db/sql"
	"techtrainingcamp-group3/pkg/db/sql/sqlAPI"
	"techtrainingcamp-group3/pkg/logger"
	"time"
)

type cb func(msg *sarama.ConsumerMessage) error

func consumeOpenEnvelope(msg *sarama.ConsumerMessage) error {
	envelope := dbmodels.Envelope{}
	err := envelope.UnmarshalBinary(msg.Value)
	if err != nil{
		logger.Sugar.Errorw("consumeOpenEnvelope message unmarshal error")
		return err
	}

	//try to update user amount and envelope status in sql
	err = sqlAPI.UpdateEnvelopeOpen(&envelope)
	if err != nil{
		return err
	}

	// If data success flush envelope to redis
	envelope.Opened = true
	if err := redisAPI.SetEnvelopeByEID(&envelope, 300*time.Second); err != nil {
		logger.Sugar.Errorw("Redis set envelop opened error", "envelope_id", envelope.EnvelopeId, "uid", envelope.Uid)
	}
	return nil
}

func consumeAddUser(msg *sarama.ConsumerMessage) error {
	tx := sql.DB.Begin()
	user := dbmodels.User{}
	user.UnmarshalBinary(msg.Value)
	if err := sql.DB.Table(
		dbmodels.User{}.TableName()).Create(
		&user).Error; err != nil {
		logger.Sugar.Debugw("CreateUserByUID", "error", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func consumeAddEnvelopeToUser(msg *sarama.ConsumerMessage) error {
	tmp, _ := strconv.Atoi(string(msg.Key))
	uid := dbmodels.UID(tmp)
	envelope := dbmodels.Envelope{}
	err := envelope.UnmarshalBinary(msg.Value)
	if err != nil{
		logger.Sugar.Errorw("ConsumeAddEnvelopeToUser UnmarshalBinary error","uid",uid)
		return err
	}

	logger.Sugar.Debugw("Consumer: AddEnvelopeToUser", "uid", uid, "envelope", envelope)
	err = sqlAPI.AddEnvelopeToUserByUID(uid,envelope)
	if err != nil{
		logger.Sugar.Errorw("AddEnvelopeToUserByUID SQL error","error",err)
		return err
	}

	//Update envelope's information in redis
	err = redisAPI.SetEnvelopeByEID(&envelope, 300*time.Second)
	if err != nil {
		logger.Sugar.Debugw("snatch", "redis set error", err, "envelope", envelope)
	}
	return nil
}

func loopConsumer(consumer sarama.Consumer, topic string, partition int, f cb) {
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		logger.Sugar.Fatalw("kafka failed to start consumer for partition",
			"err", err,
		)
		return
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		f(msg)

	}
}

func handler(topic string) func(msg *sarama.ConsumerMessage) error {
	switch topic {
	case "OpenEnvelope":
		return consumeOpenEnvelope
	case "AddUser":
		return consumeAddUser
	case "AddEnvelopeToUser":
		return consumeAddEnvelopeToUser
	default:
		return nil
	}
}
