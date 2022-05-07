package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	auth_accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/authaccounts"
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
// @Param        ecu  body      CreateTaskRequest
// @Success      201  {object}  CreateTaskResponse
// @Router       /api/v1/auth/accounts [post]
func (s *Service) CreateAccount(ctx *gin.Context) {

	var request CreateAccountRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.authAccountsRepo.Create((*auth_accounts_repo.AuthAccountCreateRow)(&request)); err != nil {
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

	msg := AccountCreatedMessage{
		PublicId: account.PublicId,
		Email:    account.Email,
		Fullname: account.Fullname,
		Position: account.Position,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	baseMsg := BaseKafkaMessage{
		Type: accountCreatedMsgType,
		Data: msgBytes,
	}
	baseMsgBytes, err := json.Marshal(baseMsg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	kafkaMsg := prepareMessage(accountTopic, strconv.Itoa(account.Id), baseMsgBytes)
	_, _, err = s.kafkaProducer.SendMessage(kafkaMsg)
	if err != nil {
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
