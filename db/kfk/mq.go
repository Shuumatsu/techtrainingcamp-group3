package kfk

import (
	"github.com/Shopify/sarama"
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/logger"
)

func init() {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{}
	msg.Topic = "test"
	msg.Value = sarama.StringEncoder("this is a good test, my message is good")

	client, err := sarama.NewSyncProducer([]string{config.Env.KafkaHost + ":" + config.Env.KafkaPort}, cfg)
	if err != nil {
		logger.Sugar.Fatalw("kafka init test producer close",
			"kafka config", cfg,
			"err", err,
		)
		return
	}

	defer client.Close()

	_, _, err = client.SendMessage(msg)
	if err != nil {
		logger.Sugar.Fatalw("kafka init test send message failed",
			"kafka config", cfg,
			"err", err,
		)
		return
	}
	logger.Sugar.Debugw("kafka init success", cfg)
}
