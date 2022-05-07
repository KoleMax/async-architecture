package tasks

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShuffleTasksResponse struct{}

// CreateTask	 godoc
// @Summary      Shuffle tasks
// @Description  Shuffle tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      201  {object}  ShuffleTasksResponse
// @Router       /api/v1/tasks/shuffle [post]
// @Security     OAuth2Password
func (s *Service) ShuffleTasks(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{AdminPosition, ManagerPosition})
	if authAccount == nil {
		return
	}

	workers, err := s.accountsRepo.ListWorkers()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(workers) == 0 {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "no workers found"})
		return
	}

	tasks, err := s.tasksRepo.ListActive()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, task := range tasks {
		randomWorkerIndex := rand.Intn(len(workers))

		if err := s.tasksRepo.Assigne(task.Id, workers[randomWorkerIndex].Id); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msg := TaskAssignedMessage{
			Id:              task.Id,
			AssignePublicId: workers[randomWorkerIndex].PublicId,
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		baseMsg := BaseKafkaMessage{
			Type: taskAssignedType,
			Data: msgBytes,
		}
		baseMsgBytes, err := json.Marshal(baseMsg)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		kafkaMsg := prepareMessage(tasksBeTopic, strconv.Itoa(task.Id), baseMsgBytes)
		_, _, err = s.producer.SendMessage(kafkaMsg)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("kafka send error: %v", err)})
			return
		}

	}

	ctx.AbortWithStatusJSON(http.StatusCreated, nil)
}
