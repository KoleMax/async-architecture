package accounting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetMyBalanceResponse struct {
	Balance int `json:"balance"`
}

// CreateTask	 godoc
// @Summary      Get my balance
// @Description  Get my balance
// @Tags         accounting
// @Accept       json
// @Produce      json
// @Success      201  {object}  GetMyBalanceResponse
// @Router       /api/v1/my-balance [get]
// @Security     OAuth2Password
func (s *Service) GetMyBalance(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{WorkerPosition})
	if authAccount == nil {
		return
	}

	accoutingAccount, err := s.accountsRepo.GetByPublicId(authAccount.PublicId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, GetMyBalanceResponse{
		Balance: accoutingAccount.Balance,
	})
}
