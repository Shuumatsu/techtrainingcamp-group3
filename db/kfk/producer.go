package kfk

import (
	"github.com/Shopify/sarama"
	"techtrainingcamp-group3/db/dbmodels"
	"strconv"
)


func OpenEnvelope(uid dbmodels.UID, envelope dbmodels.Envelope) {
	msg := &sarama.ProducerMessage{}
	data, _ := envelope.MarshalBinary()
	msg.Key = sarama.StringEncoder(strconv.FormatUint(uint64(uid), 10))
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

func AddUser(user dbmodels.User) {
	msg := &sarama.ProducerMessage{}
	data, _ := user.MarshalBinary()
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

func AddEnvelopeToUser(uid dbmodels.UID, envelope dbmodels.Envelope) {
	msg := &sarama.ProducerMessage{}
	data, _ := envelope.MarshalBinary()
	msg.Key = sarama.StringEncoder(strconv.FormatUint(uint64(uid), 10))
	msg.Value = sarama.ByteEncoder(data)
	Producer.Input() <- msg
}

var Producer sarama.AsyncProducer