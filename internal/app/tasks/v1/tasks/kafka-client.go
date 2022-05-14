package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

var brokers = []string{"localhost:29092"}

const (
	accountTopic = "accounts-stream"
	tasksTopic   = "tasks-stream"
	tasksBeTopic = "tasks"

	consumerClientId = "tasks"
)

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

func prepareMessage(topic, key string, message []byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(message),
	}

	return msg
}

func newConsumer() (sarama.Consumer, error) {
	conusemrConfig := sarama.NewConfig()
	conusemrConfig.ClientID = consumerClientId
	conusemrConfig.Consumer.Return.Errors = true
	return sarama.NewConsumer(brokers, conusemrConfig)
}

func (s *Service) consumeAccounts() error {
	partitions, _ := s.consumer.Partitions(accountTopic)
	consumer, err := s.consumer.ConsumePartition(accountTopic, partitions[0], sarama.OffsetOldest)

	if err != nil {
		fmt.Printf("Topic %v Partitions: %v", accountTopic, partitions)
		return err
	}
	fmt.Println(" Start consuming topic ", accountTopic)
	go func(topic string, consumer sarama.PartitionConsumer) {
		for {
			select {
			case consumerError := <-consumer.Errors():
				fmt.Println("consumerError: ", consumerError.Err)
				panic(consumerError.Err)
			case msg := <-consumer.Messages():
				if err := s.handleKafkaMsg(msg); err != nil {
					fmt.Println("handleKafkaMsg: ", err)
				}
			}
		}
	}(accountTopic, consumer)

	return nil
}

func (s *Service) handleKafkaMsg(msg *sarama.ConsumerMessage) error {
	var baseMsg BaseKafkaMessage
	if err := json.Unmarshal(msg.Value, &baseMsg); err != nil {
		return err
	}

	if baseMsg.Type == accountCreatedMsgType {
		var msg AccountCreatedMessage
		if err := json.Unmarshal(baseMsg.Data, &msg); err != nil {
			return err
		}
		_, err := s.accountsRepo.Create(msg.PublicId, msg.Fullname, msg.Position)
		if err != nil {
			return err
		}
	}

	return nil
}
