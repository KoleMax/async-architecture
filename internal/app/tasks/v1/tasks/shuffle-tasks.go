package tasks

import (
	"fmt"
	"math/rand"
	"net/http"

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

	tasksWithAssignePublicId := make([]taskToAssignePublicId, 0, len(tasks))

	for _, task := range tasks {
		randomWorkerIndex := rand.Intn(len(workers))

		if err := s.tasksRepo.Assigne(task.Id, workers[randomWorkerIndex].Id); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tasksWithAssignePublicId = append(tasksWithAssignePublicId, taskToAssignePublicId{
			Task:            &task,
			AssignePublicId: workers[randomWorkerIndex].PublicId,
		})
	}

	if err := s.sendMultipleTaskAssignedV1(tasksWithAssignePublicId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("kafka sender error: %v", err.Error())})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusCreated, nil)
}
