package dbrepo

import (
	"database/sql"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

// The following blocks of code is used for testing purposes
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
