package accounting

import (
	"net/http"

	transactions_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/accounting/transactions"
	"github.com/gin-gonic/gin"
)

type ListMyTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// CreateTask	 godoc
// @Summary      List my transactions
// @Description  List my transactions
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Success      201  {object}  ListMyTransactionsResponse
// @Router       /api/v1/transactions/my [get]
// @Security     OAuth2Password
func (s *Service) ListMyTransactions(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{WorkerPosition})
	if authAccount == nil {
		return
	}

	accountingAccount, err := s.accountsRepo.GetByPublicId(authAccount.PublicId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactions, err := s.transactionsRepo.GetByAccountId(accountingAccount.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		transactionType := TransactionTypeCredit
		cost := transaction.CostDone
		if transaction.Type == transactions_repo.TransactionTypeDebit {
			transactionType = TransactionTypeWriteOff
			cost = transaction.CostAssigne
		}

		result = append(result, Transaction{
			Id:              transaction.Id,
			TaskId:          transaction.TaskId,
			TaskDescription: transaction.TaskDescription,
			EventCreatedAt:  transaction.EventCreatedAt,
			Type:            transactionType,
			Cost:            cost,
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, ListMyTransactionsResponse{
		Transactions: result,
	})
}
