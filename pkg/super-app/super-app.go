package superapp

import (
	"github.com/KoleMax/async-architecture/pkg/super-app/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Service interface {
	SetupRoutes(gin.IRouter)
}

type SuperApp struct {
	router *gin.Engine
}

func New(services []Service) (*SuperApp, error) {
	router := gin.Default()

	setupSwagger(router)
	setupMiddleware(router)

	for _, s := range services {
		s.SetupRoutes(router)
	}

	return &SuperApp{
		router,
	}, nil
}

func setupMiddleware(router *gin.Engine) {
	router.Use(gin.CustomRecovery(middleware.Recover))
	router.Use(middleware.LogLevelOverride)
	// router.Use(middleware.ParseTokenMiddleware)
}

func setupSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
