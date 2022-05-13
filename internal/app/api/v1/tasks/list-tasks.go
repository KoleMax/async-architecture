package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// CreateTask	 godoc
// @Summary      List all tasks
// @Description  List all tasks
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      201  {object}  ListTasksResponse
// @Router       /api/v1/tasks [get]
// @Security     OAuth2Password
func (s *Service) ListTasks(ctx *gin.Context) {
	authAccount := Authorize(ctx, []string{AdminPosition, ManagerPosition})
	if authAccount == nil {
		return
	}

	tasks, err := s.tasksRepo.ListAll()
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

	ctx.AbortWithStatusJSON(http.StatusOK, ListTasksResponse{
		Tasks: result,
	})
}
