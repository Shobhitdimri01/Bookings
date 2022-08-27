package dbrepo

import (
	"database/sql"

	"github.com/Shobhitdimri01/Bookings/internal/config"
	"github.com/Shobhitdimri01/Bookings/internal/repository"
)

//This struct and function give ease to connect to any database(MongoDb,MariaDb ...etc) eg:-->Here we are Connecting to Postgres
type postgresDBRepo struct {
	App *config.AppConfig
	Db 	*sql.DB	
}

func NewPostgresRepo(conn *sql.DB,a *config.AppConfig) repository.DataBaseRepo{
	return &postgresDBRepo{
		App:a,
		Db:conn,
	}
}