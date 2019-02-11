package database

import (
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/go-pg/pg"
)

func Connect(env *core.Environment) *pg.DB {
	options := &pg.Options{
		User:     env.DbUser,
		Password: env.DbPassword,
		Database: env.DbName,
	}

	db := pg.Connect(options)
	if db == nil {
		panic("Could not connect to database")
	}

	return db
}

func Close(db *pg.DB) {
	dbCloseError := db.Close()
	if dbCloseError != nil {
		panic("Could not close connection to database")
	}
}
