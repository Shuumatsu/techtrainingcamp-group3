package kfk

import (
	"context"
	"strconv"
	"techtrainingcamp-group3/pkg/db/dbmodels"
	"techtrainingcamp-group3/pkg/logger"
	"time"

	"github.com/Shopify/sarama"
)

func OpenEnvelope(uid dbmodels.UID, envelope *dbmodels.Envelope) error{
	msg := &sarama.ProducerMessage{}
	msg.Topic = "OpenEnvelope"
	data, _ := envelope.MarshalBinary()
	msg.Key = sarama.StringEncoder(strconv.FormatUint(uint64(uid), 10))
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg

	//set 30 seconds' time out
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx,30*time.Second)
	defer cancel()

	select {
	//put message into kafka successfully
	case <-Producer.Successes():
		return nil
	case err := <-Producer.Errors():
		return  err
		// time out
	case <-ctx.Done():
		return ctx.Err()
	}
}

func AddUser(user dbmodels.User) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = "AddUser"
	data, _ := user.MarshalBinary()
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

func AddEnvelopeToUser(uid dbmodels.UID, envelope dbmodels.Envelope) error{
	msg := &sarama.ProducerMessage{}
	msg.Topic = "AddEnvelopeToUser"
	data, err := envelope.MarshalBinary()
	if err != nil{
		logger.Sugar.Errorw("kafka producer error:","uid",uid)
		return err
	}

	msg.Key = sarama.StringEncoder(strconv.FormatUint(uint64(uid), 10))
	msg.Value = sarama.ByteEncoder(data)

	Producer.Input() <- msg

	//set 30 seconds' time out
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx,30*time.Second)
	defer cancel()

	select {
	case <-Producer.Successes():
		return nil
	case err := <-Producer.Errors():
		return  err
	case <-ctx.Done():
		return ctx.Err()
	}
}

var Producer sarama.AsyncProducer
