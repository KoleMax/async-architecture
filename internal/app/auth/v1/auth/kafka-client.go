package auth

import (
	"strconv"

	"github.com/Shopify/sarama"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/auth/accounts"
	auth_kafka_msgs "github.com/KoleMax/async-architecture/pkg/kafka-schemas/auth"
)

var brokers = []string{"localhost:29092"}

const (
	senderName = "auth"

	accountTopic = "accounts-stream"
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

func (s *Service) sendAccountCreatedV1(account accounts_repo.AccountGetRow) error {
	msg := auth_kafka_msgs.AccountCreatedV1{
		PublicId: account.PublicId,
		Email:    account.Email,
		Fullname: account.Fullname,
		Position: account.Position,
	}
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	baseMsg := auth_kafka_msgs.Base{
		Version: "v1",
		Sender:  senderName,
		Type:    auth_kafka_msgs.MessageType_AccountCreated,
		Sended:  timestamppb.Now(),
		Data:    data,
	}
	baseData, err := proto.Marshal(&baseMsg)
	if err != nil {
		return err
	}

	kafkaMsg := prepareMessage(accountTopic, strconv.Itoa(account.Id), baseData)
	_, _, err = s.kafkaProducer.SendMessage(kafkaMsg)
	return err
}
