package accounting

import (
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"

	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounting/accounts"
	payments_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounting/payments"
	tasks_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounting/tasks"
	transactions_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounting/transactions"
)

type Service struct {
	accountsRepo     *accounts_repo.Repository
	tasksRepo        *tasks_repo.Repository
	transactionsRepo *transactions_repo.Repository
	paymentsRepo     *payments_repo.Repository
	consumer         sarama.Consumer
}

func New(
	accountsRepo *accounts_repo.Repository,
	tasksRepo *tasks_repo.Repository,
	transactionsRepo *transactions_repo.Repository,
	paymentsRepo *payments_repo.Repository,
) *Service {
	consumer, err := newConsumer()
	if err != nil {
		panic(err)
	}

	service := &Service{
		accountsRepo:     accountsRepo,
		tasksRepo:        tasksRepo,
		transactionsRepo: transactionsRepo,
		paymentsRepo:     paymentsRepo,
		consumer:         consumer,
	}

	service.consumeAccounts()
	service.consumeTasks()

	return service
}

func (s *Service) SetupRoutes(router gin.IRouter) {
	router.GET("/api/v1/my-balance", s.GetMyBalance)
	router.GET("/api/v1/transactions/my", s.ListMyTransactions)
	router.GET("/api/v1/payments/my", s.ListMyPayments)
}
