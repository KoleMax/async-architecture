package main

import (
	"context"

	docs "github.com/KoleMax/async-architecture/docs/tasks" // TODO: configure swagger
	tasks_service "github.com/KoleMax/async-architecture/internal/app/tasks/v1/tasks"
	"github.com/KoleMax/async-architecture/internal/pkg/config"
	"github.com/KoleMax/async-architecture/internal/pkg/db"
	accounts_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/tasks/accounts"
	tasks_repo "github.com/KoleMax/async-architecture/internal/pkg/repository/tasks/tasks"
	"github.com/KoleMax/async-architecture/pkg/logging"
	superapp "github.com/KoleMax/async-architecture/pkg/super-app"
)

func main() {
	ctx := context.Background()

	docs.Init(config.Get().Auth.IssuerUrl)

	if err := logging.SetLevel(logging.GetGlobalLogger(), config.Get().Logging.Level); err != nil {
		logging.Panicf(ctx, "logging.SetLevel: %v", err)
	}

	dbConnection, err := db.New()
	if err != nil {
		logging.Panicf(ctx, "db.New: %v", err)
	}

	tasksRepo := tasks_repo.New(dbConnection)
	accountsRepo := accounts_repo.New(dbConnection)

	var services []superapp.Service
	services = append(services,
		tasks_service.New(tasksRepo, accountsRepo),
	)

	app, err := superapp.New(services)
	if err != nil {
		logging.Panicf(ctx, "superapp.New: %v", err)
	}

	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if err := app.Run(); err != nil {
		logging.Panicf(ctx, "app.Run: %v", err)
	}
}
