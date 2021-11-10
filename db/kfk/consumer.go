package kfk

import (
	"fmt"
	"strconv"
	"techtrainingcamp-group3/db/dbmodels"
	"techtrainingcamp-group3/db/sql"
	"techtrainingcamp-group3/logger"
	"github.com/Shopify/sarama"
	"gorm.io/gorm"
)

type cb func(msg *sarama.ConsumerMessage) error

func consumeOpenEnvelope(msg *sarama.ConsumerMessage) error {
	tx := sql.DB.Begin()
	tmp, _ := strconv.Atoi(string(msg.Key))
	uid := dbmodels.UID(tmp)
	envelope := dbmodels.Envelope{}
	envelope.UnmarshalBinary(msg.Value)
	logger.Sugar.Debugw("Consumer: OpenEnvelopeByEID", "uid", uid, "envelope", envelope)
	if err := tx.Model(
		&envelope).Update("opened", true).Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return err
	}
	if err := tx.Table(
		dbmodels.User{}.TableName()).Where("uid", uid).Update(
		"amount", gorm.Expr(
			"amount + ?", envelope.Value)).Error; err != nil {
		logger.Sugar.Debugw("OpenEnvelopeByEID", "error", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
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
	envelope.UnmarshalBinary(msg.Value)
	logger.Sugar.Debugw("Consumer: AddEnvelopeToUser", "uid", uid, "envelope", envelope)
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

func handler(topic string) func(msg *sarama.ConsumerMessage) error{
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
