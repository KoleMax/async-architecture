package accounting

import (
	"github.com/gin-gonic/gin"
)

type Service struct {
}

func New() *Service {

	service := &Service{}

	return service
}

func (s *Service) SetupRoutes(router gin.IRouter) {}
