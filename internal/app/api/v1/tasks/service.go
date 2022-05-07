package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"

	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounts"
	tasks_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/tasks"
)

var brokers = []string{"localhost:9092"}

const (
	accountTopic = "accounts-stream"
	tasksTopic   = "tasks-stream"
	tasksBeTopic = "tasks"
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

type Service struct {
	tasksRepo    *tasks_repo.Repository
	accountsRepo *accounts_repo.Repository
	consumer     sarama.Consumer
	producer     sarama.SyncProducer
}

func New(tasksRepo *tasks_repo.Repository, accountsRepo *accounts_repo.Repository) *Service {
	conusemrConfig := sarama.NewConfig()
	conusemrConfig.ClientID = "tasks"
	conusemrConfig.Consumer.Return.Errors = true

	// Create new consumer
	consumer, err := sarama.NewConsumer(brokers, conusemrConfig)
	if err != nil {
		panic(err)
	}

	producer, err := newProducer()
	if err != nil {
		panic(err)
	}

	service := &Service{
		tasksRepo:    tasksRepo,
		accountsRepo: accountsRepo,
		consumer:     consumer,
		producer:     producer,
	}

	service.consumeAccounts()

	return service
}

func (s *Service) SetupRoutes(router gin.IRouter) {
	router.GET("/api/v1/tasks/", s.ListTasks)
	router.GET("/api/v1/tasks/my", s.ListMyTasks)

	router.POST("/api/v1/tasks", s.CreateTask)

	router.POST("/api/v1/tasks/:id/complete", s.CompleteTask)
	router.POST("/api/v1/tasks/shuffle", s.ShuffleTasks)
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
	if err := json.Unmarshal(msg.Value, baseMsg); err != nil {
		return err
	}

	if baseMsg.Type == accountCreatedMsgType {
		var msg AccountCreatedMessage
		if err := json.Unmarshal(baseMsg.Data, msg); err != nil {
			return err
		}
		_, err := s.accountsRepo.Create(msg.PublicId, msg.Fullname, msg.Position)
		if err != nil {
			return err
		}
	}

	return nil
}
