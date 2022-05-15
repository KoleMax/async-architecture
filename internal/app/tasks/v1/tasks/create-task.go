package tasks

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	JiraId      string `json:"jira_id"`
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

	task, err := s.tasksRepo.Create(workers[randomWorkerIndex].Id, request.Title, request.JiraId, request.Description)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.sendTaskAddedV1(task, authAccount.PublicId); err != nil {
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
