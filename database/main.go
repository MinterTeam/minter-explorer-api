package database

import (
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/go-pg/pg"
	"os"
)

func Connect(env *core.Environment) *pg.DB {
	options := &pg.Options{
		User:     env.DbUser,
		Password: env.DbPassword,
		Database: env.DbName,
	}

	db := pg.Connect(options)
	if db == nil {
		println("Could not connect to database")
		os.Exit(100)
	}

	return db
}

func Close(db *pg.DB) {
	dbCloseError := db.Close()
	if dbCloseError != nil {
		println("Could not close connection to database")
		os.Exit(100)
	}
}
