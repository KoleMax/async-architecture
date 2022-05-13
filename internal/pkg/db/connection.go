package db

import (
	"fmt"

	"github.com/KoleMax/async-architecture/internal/pkg/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func New() (*sqlx.DB, error) {
	cfg := config.Get().Db
	pgxCfg, err := pgx.ParseConfig(
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.Name,
			cfg.User,
			cfg.Password,
		),
	)
	if err != nil {
		return nil, err
	}

	pgxCfg.LogLevel = pgx.LogLevelError
	db := stdlib.OpenDB(*pgxCfg)
	db.SetMaxOpenConns(cfg.MaxConnsNum)
	db.SetMaxIdleConns(cfg.MinConnsNum)

	return sqlx.NewDb(db, "pgx"), nil
}
