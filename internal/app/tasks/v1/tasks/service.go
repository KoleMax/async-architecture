package tasks

import (
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"

	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/tasks/accounts"
	tasks_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/tasks/tasks"
)

type Service struct {
	tasksRepo    *tasks_repo.Repository
	accountsRepo *accounts_repo.Repository
	consumer     sarama.Consumer
	producer     sarama.SyncProducer
}

func New(tasksRepo *tasks_repo.Repository, accountsRepo *accounts_repo.Repository) *Service {
	consumer, err := newConsumer()
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
