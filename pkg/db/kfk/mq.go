package kfk

import (
	"techtrainingcamp-group3/config"
	"techtrainingcamp-group3/pkg/logger"

	"github.com/Shopify/sarama"
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
			"err", err,
		)
		return
	}

	_, _, err = client.SendMessage(msg)
	if err != nil {
		logger.Sugar.Fatalw("kafka init test send message failed",
			"err", err,
		)
		return
	}
	Producer, _ = sarama.NewAsyncProducer([]string{config.Env.KafkaHost + ":" + config.Env.KafkaPort}, cfg)
	for _, topic := range config.Env.KafkaTopics {
		consumer, err := sarama.NewConsumer([]string{config.Env.KafkaHost + ":" + config.Env.KafkaPort}, cfg)

		if err != nil {
			logger.Sugar.Fatalw("kafka fail to start consumer", "err", err)
			return
		}
		partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区

		if err != nil {
			logger.Sugar.Fatalw("kafka fail to get list of partition", "err", err)
			return
		}

		for partition := range partitionList { // 遍历所有的分区
			// 针对每个分区创建一个对应的分区消费者
			go loopConsumer(consumer, topic, partition, handler(topic))
		}
	}
	logger.Sugar.Debugw("kafka init success", cfg)
}
