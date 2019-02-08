package main

import (
	"flag"
	"github.com/MinterTeam/minter-explorer-api/api"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/database"
	"os"
)

func main() {
	env := initEnvironment()
	db := database.Connect(env)
	defer database.Close(db)

	api.Run(db)
}

func initEnvironment() *core.Environment {
	env := &core.Environment{
		DbName:     *flag.String("db_name", "", "DB name"),
		DbUser:     *flag.String("db_user", "", "DB user"),
		DbPassword: *flag.String("db_password", "", "DB password"),
	}
	flag.Parse()

	if env.DbUser == `` {
		dbUser := os.Getenv("EXPLORER_DB_USER")
		env.DbUser = dbUser
	}
	if env.DbName == `` {
		dbName := os.Getenv("EXPLORER_DB_NAME")
		env.DbName = dbName
	}
	if env.DbPassword == `` {
		dbPassword := os.Getenv("EXPLORER_DB_PASSWORD")
		env.DbPassword = dbPassword
	}

	return env
}
