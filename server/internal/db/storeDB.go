package db

import (
	"database/sql"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/logg"
	_ "github.com/lib/pq"
)

// StoreDB database.
type StoreDB struct {
	Logger logg.Logger
	DB     *sql.DB
}

func NewStoreDB(config config.Config, logger logg.Logger) *StoreDB {
	db, err := sql.Open("postgres", config.GetConnDB())
	if err != nil {
		logger.RaiseFatal(err, "error open database", nil)
	}

	return &StoreDB{
		Logger: logger,
		DB:     db,
	}
}

func (sd *StoreDB) CloseDB() {
	sd.DB.Close()
}

func (sd *StoreDB) GetDB() *sql.DB {
	return sd.DB
}

func (sd *StoreDB) CheckOpen() bool {
	err := sd.DB.Ping()
	if err != nil {
		return false
	}
	return true
}
