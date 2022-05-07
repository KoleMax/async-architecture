package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/Shopify/sarama"

	auth_accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/authaccounts"
)

var brokers = []string{"localhost:9092"}

const accountTopic = "accounts-stream"

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

type Service struct {
	authAccountsRepo *auth_accounts_repo.Repository
	kafkaProducer    sarama.SyncProducer
}

func New(authAccountsRepo *auth_accounts_repo.Repository) (*Service, error) {
	kafkaProdcuer, err := newProducer()
	if err != nil {
		return nil, err
	}
	return &Service{
		authAccountsRepo: authAccountsRepo,
		kafkaProducer:    kafkaProdcuer,
	}, nil
}

func (s *Service) SetupRoutes(router gin.IRouter) {
	router.POST("/api/v1/auth/authenticate", s.Authenticate)

	router.POST("/api/v1/auth/accounts", s.CreateAccount)
}
