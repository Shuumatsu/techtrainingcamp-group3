package kfk

import (
	"strconv"
	"techtrainingcamp-group3/pkg/db/dbmodels"

	"github.com/Shopify/sarama"
)

func OpenEnvelope(uid dbmodels.UID, envelope dbmodels.Envelope) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = "OpenEnvelope"
	data, _ := envelope.MarshalBinary()
	msg.Key = sarama.StringEncoder(strconv.FormatUint(uint64(uid), 10))
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

func AddUser(user dbmodels.User) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = "AddUser"
	data, _ := user.MarshalBinary()
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

func AddEnvelopeToUser(uid dbmodels.UID, envelope dbmodels.Envelope) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = "AddEnvelopeToUser"
	data, _ := envelope.MarshalBinary()
	msg.Key = sarama.StringEncoder(strconv.FormatUint(uint64(uid), 10))
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

var Producer sarama.AsyncProducer