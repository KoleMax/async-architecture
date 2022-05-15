package tasks

import (
	"fmt"
	"strconv"

	"github.com/Shopify/sarama"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	tasks_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/tasks/tasks"
	auth_kafka_msgs "github.com/KoleMax/async-architecture/pkg/kafka-schemas/auth"
	task_kafka_msgs "github.com/KoleMax/async-architecture/pkg/kafka-schemas/tasks"
)

var brokers = []string{"localhost:29092"}

const (
	senderName = "tasks"

	accountTopic = "accounts-stream"
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
				if err := s.handleKafkaAccountsMsg(msg); err != nil {
					fmt.Println("handleKafkaMsg: ", err)
				}
			}
		}
	}(accountTopic, consumer)

	return nil
}

func (s *Service) handleKafkaAccountsMsg(msg *sarama.ConsumerMessage) error {
	var baseMsg auth_kafka_msgs.Base
	if err := proto.Unmarshal(msg.Value, &baseMsg); err != nil {
		return err
	}

	if baseMsg.Type == auth_kafka_msgs.MessageType_AccountCreated && baseMsg.Version == "v1" {
		var msg auth_kafka_msgs.AccountCreatedV1
		if err := proto.Unmarshal(baseMsg.Data, &msg); err != nil {
			return err
		}
		_, err := s.accountsRepo.Create(msg.PublicId, msg.Fullname, msg.Position)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sendTaskAddedV1(task *tasks_repo.Task, assignePublicId string) error {
	msg := task_kafka_msgs.TaskAddedV2{
		TaskPublicId:    task.PublicId,
		AssignePublicId: assignePublicId,
		Title:           task.Title,
		JiraId:          task.JiraId,
		Description:     task.Description,
	}
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	baseMsg := task_kafka_msgs.Base{
		Version: "v2",
		Sender:  senderName,
		Type:    task_kafka_msgs.MessageType_TaskAdded,
		Sended:  timestamppb.Now(),
		Data:    data,
	}
	baseData, err := proto.Marshal(&baseMsg)
	if err != nil {
		return err
	}

	kafkaMsg := prepareMessage(tasksBeTopic, strconv.Itoa(task.Id), baseData)
	_, _, err = s.producer.SendMessage(kafkaMsg)
	return err
}

type taskToAssignePublicId struct {
	Task            *tasks_repo.Task
	AssignePublicId string
}

func (s *Service) sendMultipleTaskAssignedV1(tasks []taskToAssignePublicId) error {
	msgs := make([]*sarama.ProducerMessage, 0, len(tasks))

	for _, task := range tasks {
		msg := task_kafka_msgs.TaskAssignedV1{
			TaskPublicId:    task.Task.PublicId,
			AssignePublicId: task.AssignePublicId,
		}
		data, err := proto.Marshal(&msg)
		if err != nil {
			return err
		}

		baseMsg := task_kafka_msgs.Base{
			Version: "v1",
			Sender:  senderName,
			Type:    task_kafka_msgs.MessageType_TaskAssigned,
			Sended:  timestamppb.Now(),
			Data:    data,
		}
		baseData, err := proto.Marshal(&baseMsg)
		if err != nil {
			return err
		}

		kafkaMsg := prepareMessage(tasksBeTopic, strconv.Itoa(task.Task.Id), baseData)
		msgs = append(msgs, kafkaMsg)
	}

	return s.producer.SendMessages(msgs)
}

func (s *Service) sendTaskCompletedV1(task *tasks_repo.Task, assignePublicId string) error {
	msg := task_kafka_msgs.TaskCompletedV1{
		TaskPublicId:    task.PublicId,
		AssignePublicId: assignePublicId,
	}
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	baseMsg := task_kafka_msgs.Base{
		Version: "v1",
		Sender:  senderName,
		Type:    task_kafka_msgs.MessageType_TaskCompleted,
		Sended:  timestamppb.Now(),
		Data:    data,
	}
	baseData, err := proto.Marshal(&baseMsg)
	if err != nil {
		return err
	}

	kafkaMsg := prepareMessage(tasksBeTopic, strconv.Itoa(task.Id), baseData)
	_, _, err = s.producer.SendMessage(kafkaMsg)
	return err
}
