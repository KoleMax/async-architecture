package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListMyTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// CreateTask	 godoc
// @Summary      List my tasks
// @Description  List my tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      201  {object}  ListMyTasksResponse
// @Router       /api/v1/tasks/my [get]
// @Security     OAuth2Password
func (s *Service) ListMyTasks(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{WorkerPosition})
	if authAccount == nil {
		return
	}

	tasks, err := s.tasksRepo.List(authAccount.PublicId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, Task{
			Id:          task.Id,
			AssigneId:   task.AssigneId,
			Description: task.Description,
			Status:      task.Status,
		})
	}

	ctx.AbortWithStatusJSON(http.StatusOK, ListMyTasksResponse{
		Tasks: result,
	})
}
