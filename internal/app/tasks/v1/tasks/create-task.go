package tasks

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	Task Task `json:"task"`
}

// CreateTask	 godoc
// @Summary      Create new task
// @Description  Create new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        ecu  body      CreateTaskRequest  true  "Add task"
// @Success      201  {object}  CreateTaskResponse
// @Router       /api/v1/tasks [post]
// @Security     OAuth2Password
func (s *Service) CreateTask(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{AdminPosition, ManagerPosition, AccountantPosition, WorkerPosition})
	if authAccount == nil {
		return
	}

	var request CreateTaskRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	randomWorkerIndex := rand.Intn(len(workers))

	task, err := s.tasksRepo.Create(workers[randomWorkerIndex].Id, request.Description)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := TaskAddedMessage{
		Id:              task.Id,
		AssignePublicId: workers[randomWorkerIndex].PublicId,
		Description:     task.Description,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	baseMsg := BaseKafkaMessage{
		Type: taskAddedType,
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

	ctx.AbortWithStatusJSON(http.StatusCreated, &CreateTaskResponse{
		Task: Task{
			Id:          task.Id,
			AssigneId:   task.AssigneId,
			Description: task.Description,
			Status:      task.Status,
		},
	})
}
