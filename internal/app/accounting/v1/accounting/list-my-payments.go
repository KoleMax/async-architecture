package accounting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListMyPaymentsResponse struct {
	Payments []Payment `json:"payments"`
}

// CreateTask	 godoc
// @Summary      List my payments
// @Description  List my payments
// @Tags         payments
// @Accept       json
// @Produce      json
// @Success      201  {object}  ListMyPaymentsResponse
// @Router       /api/v1/payments/my [get]
// @Security     OAuth2Password
func (s *Service) ListMyPayments(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{WorkerPosition})
	if authAccount == nil {
		return
	}

	accountingAccount, err := s.accountsRepo.GetByPublicId(authAccount.PublicId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payments, err := s.paymentsRepo.GetByAccountId(accountingAccount.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]Payment, 0, len(payments))
	for _, payment := range payments {

		result = append(result, Payment{
			Id:             payment.Id,
			Amount:         payment.Amount,
			Status:         payment.Status,
			CreatedAt:      payment.CreatedAt,
			EventCreatedAt: payment.EventCreatedAt,
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, ListMyPaymentsResponse{
		Payments: result,
	})
}
