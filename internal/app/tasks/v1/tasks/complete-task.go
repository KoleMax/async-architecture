package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CompleteTask	 godoc
// @Summary      Complete task
// @Description  Complete task
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

	task, err := s.tasksRepo.GetById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.sendTaskCompletedV1(task, authAccount.PublicId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("kafka send error: %v", err)})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}
