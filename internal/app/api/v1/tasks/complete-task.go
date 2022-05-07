package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CompleteTask	 godoc
// @Summary      Complete firmware
// @Description  Complete firmware
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Task Id"
// @Success      201
// @Router       /api/v1/tasks/{id}/complete [post]
// @Security     OAuth2Password
func (s *Service) CompleteTask(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{AdminPosition, ManagerPosition})
	if authAccount == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = s.tasksRepo.Complete(id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := TaskCompletedMessage{
		Id:              id,
		AssignePublicId: authAccount.PublicId,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	baseMsg := BaseKafkaMessage{
		Type: taskCompletedType,
		Data: msgBytes,
	}
	baseMsgBytes, err := json.Marshal(baseMsg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	kafkaMsg := prepareMessage(tasksBeTopic, strconv.Itoa(id), baseMsgBytes)
	_, _, err = s.producer.SendMessage(kafkaMsg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("kafka send error: %v", err)})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}
