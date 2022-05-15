package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/Shopify/sarama"

	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/auth/accounts"
)

// @securitydefinitions.oauth2.application  OAuth2Application
// @tokenUrl                                http://localhost:3000/oauth2/token
// @securitydefinitions.oauth2.implicit     OAuth2Implicit
// @authorizationUrl                        http://localhost:3000/oauth/authorize
// @scope.basket-api
// @securitydefinitions.oauth2.password    OAuth2Password
// @tokenUrl                               http://localhost:3000/oauth/token
// @securitydefinitions.oauth2.accessCode  OAuth2AccessCode
// @tokenUrl                               http://localhost:3000/oauth/token
// @authorizationUrl                       http://localhost:3000/oauth/authorize
// @scope.basket-api
type Service struct {
	authAccountsRepo *accounts_repo.Repository
	kafkaProducer    sarama.SyncProducer
}

func New(authAccountsRepo *accounts_repo.Repository) (*Service, error) {
	kafkaProdcuer, err := newProducer()
	if err != nil {
		return nil, err
	}
	return &Service{
		authAccountsRepo: authAccountsRepo,
		kafkaProducer:    kafkaProdcuer,
	}, nil
}

func (s *Service) SetupRoutes(router gin.IRouter) {
	router.POST("/api/v1/auth/authenticate", s.Authenticate)

	router.POST("/api/v1/auth/accounts", s.CreateAccount)
}
