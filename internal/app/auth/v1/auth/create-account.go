package auth

import (
	"fmt"
	"net/http"

	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/auth/accounts"
	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"full_name"`
	Position string `json:"position"`
}

type CreateAccountResponse struct {
	AuthAccount AuthAccount `json:"account"`
}

// CreateTask	 godoc
// @Summary      Create new auth account
// @Description  Create new auth account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      CreateAccountRequest  true  "Add account"
// @Success      201      {object}  CreateAccountResponse
// @Router       /api/v1/auth/accounts [post]
func (s *Service) CreateAccount(ctx *gin.Context) {

	var request CreateAccountRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.authAccountsRepo.Create((*accounts_repo.AccountCreateRow)(&request)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account, err := s.authAccountsRepo.GetByEmail(request.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if account == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "no such account after creation"})
		return
	}

	if err := s.sendAccountCreatedV1(*account); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("kafka send error: %v", err)})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusCreated, CreateAccountResponse{
		AuthAccount: AuthAccount{
			PublicId: account.PublicId,
			Email:    account.Email,
			Fullname: account.Fullname,
			Position: account.Position,
		},
	})
}
