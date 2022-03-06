package dbrepo

import (
	"database/sql"
	"github.com/zubsingh/bookings/internal/config"
	"github.com/zubsingh/bookings/internal/repository"
)

// this one is created to swap db if required with any other type
// eg: sqlite define type sqliteDBRepo and Add NewSqliteRepo

type postgresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresRepo{
		App: a,
		DB:  conn,
	}
}
