package accounting

import (
	"fmt"

	"github.com/Shopify/sarama"

	"google.golang.org/protobuf/proto"

	transactions_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounting/transactions"
	auth_kafka_msgs "github.com/KoleMax/async-architecture/pkg/kafka-schemas/auth"
	task_kafka_msgs "github.com/KoleMax/async-architecture/pkg/kafka-schemas/tasks"
)

var brokers = []string{"localhost:29092"}

const (
	accountTopic = "accounts-stream"
	tasksBeTopic = "tasks"

	consumerClientId = "tasks"
)

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
		return s.accountsRepo.Create(msg.PublicId, msg.Fullname, msg.Position)
	}

	return nil
}

func (s *Service) consumeTasks() error {
	partitions, _ := s.consumer.Partitions(accountTopic)
	consumer, err := s.consumer.ConsumePartition(tasksBeTopic, partitions[0], sarama.OffsetOldest)

	if err != nil {
		fmt.Printf("Topic %v Partitions: %v", tasksBeTopic, partitions)
		return err
	}
	fmt.Println(" Start consuming topic ", tasksBeTopic)
	go func(topic string, consumer sarama.PartitionConsumer) {
		for {
			select {
			case consumerError := <-consumer.Errors():
				fmt.Println("consumerError: ", consumerError.Err)
				panic(consumerError.Err)
			case msg := <-consumer.Messages():
				if err := s.handleKafkaTasksMsg(msg); err != nil {
					fmt.Println("handleKafkaMsg: ", err)
				}
			}
		}
	}(accountTopic, consumer)

	return nil
}

func (s *Service) handleKafkaTasksMsg(msg *sarama.ConsumerMessage) error {
	var baseMsg task_kafka_msgs.Base
	if err := proto.Unmarshal(msg.Value, &baseMsg); err != nil {
		return err
	}

	if baseMsg.Type == task_kafka_msgs.MessageType_TaskAdded && baseMsg.Version == "v1" {
		var msg task_kafka_msgs.TaskAddedV1
		if err := proto.Unmarshal(baseMsg.Data, &msg); err != nil {
			return err
		}

		if err := s.tasksRepo.Create(msg.TaskPublicId, getRandomDoneCost(), getRandomAssigneCost(), msg.Title, "", msg.Description); err != nil {
			return err
		}
		
		accountingAccount, err := s.accountsRepo.GetByPublicId(msg.AssignePublicId)
		if err != nil {
			return err
		}

		task, err := s.tasksRepo.GetByPublicId(msg.TaskPublicId) 
		if err != nil {
			return err
		}

		if err := s.accountsRepo.SetBalance(accountingAccount.Id, accountingAccount.Balance - task.CostAssigne); err != nil {
			return nil
		}
		return s.transactionsRepo.Create(accountingAccount.Id, task.Id, baseMsg.Sended.AsTime(), transactions_repo.TransactionTypeDebit)
	} else if baseMsg.Type == task_kafka_msgs.MessageType_TaskAdded && baseMsg.Version == "v2" {
		var msg task_kafka_msgs.TaskAddedV2
		if err := proto.Unmarshal(baseMsg.Data, &msg); err != nil {
			return err
		}

		if err := s.tasksRepo.Create(msg.TaskPublicId, getRandomDoneCost(), getRandomAssigneCost(), msg.Title, msg.JiraId, msg.Description); err != nil {
			return err
		}
		
		accountingAccount, err := s.accountsRepo.GetByPublicId(msg.AssignePublicId)
		if err != nil {
			return err
		}

		task, err := s.tasksRepo.GetByPublicId(msg.TaskPublicId) 
		if err != nil {
			return err
		}

		if err := s.accountsRepo.SetBalance(accountingAccount.Id, accountingAccount.Balance - task.CostAssigne); err != nil {
			return nil
		}

		return s.transactionsRepo.Create(accountingAccount.Id, task.Id, baseMsg.Sended.AsTime(), transactions_repo.TransactionTypeDebit)
	} else if baseMsg.Type == task_kafka_msgs.MessageType_TaskAssigned && baseMsg.Version == "v1" {
		var msg task_kafka_msgs.TaskAssignedV1
		if err := proto.Unmarshal(baseMsg.Data, &msg); err != nil {
			return err
		}
		
		accountingAccount, err := s.accountsRepo.GetByPublicId(msg.AssignePublicId)
		if err != nil {
			return err
		}

		task, err := s.tasksRepo.GetByPublicId(msg.TaskPublicId) 
		if err != nil {
			return err
		}

		if err := s.accountsRepo.SetBalance(accountingAccount.Id, accountingAccount.Balance - task.CostAssigne); err != nil {
			return nil
		}

		return s.transactionsRepo.Create(accountingAccount.Id, task.Id, baseMsg.Sended.AsTime(), transactions_repo.TransactionTypeDebit)
	} else if baseMsg.Type == task_kafka_msgs.MessageType_TaskCompleted && baseMsg.Version == "v1" {
		var msg task_kafka_msgs.TaskCompletedV1
		if err := proto.Unmarshal(baseMsg.Data, &msg); err != nil {
			return err
		}

		accountingAccount, err := s.accountsRepo.GetByPublicId(msg.AssignePublicId)
		if err != nil {
			return err
		}

		task, err := s.tasksRepo.GetByPublicId(msg.TaskPublicId) 
		if err != nil {
			return err
		}

		if err := s.accountsRepo.SetBalance(accountingAccount.Id, accountingAccount.Balance + task.CostDone); err != nil {
			return nil
		}

		return s.transactionsRepo.Create(accountingAccount.Id, task.Id, baseMsg.Sended.AsTime(), transactions_repo.TransactionTypeCredit)
	}

	return nil
}

func getRandomDoneCost() int {
	return 0
}

func getRandomAssigneCost() int {
	return 0
}
