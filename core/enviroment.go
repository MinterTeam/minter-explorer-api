package core

import (
	"flag"
	"os"
)

type Environment struct {
	DbName     string
	DbUser     string
	DbPassword string
	BaseCoin   string
}

func NewEnvironment() *Environment {
	dbName := flag.String("db_name", "", "DB name")
	dbUser := flag.String("db_user", "", "DB user")
	dbPassword := flag.String("db_password", "", "DB password")
	baseCoin := flag.String("base_coin", "", "Base Coin")
	flag.Parse()

	env := Environment{
		DbName:     *dbName,
		DbUser:     *dbUser,
		DbPassword: *dbPassword,
		BaseCoin:   *baseCoin,
	}

	if env.DbUser == "" {
		dbUser := os.Getenv("EXPLORER_DB_USER")
		env.DbUser = dbUser
	}

	if env.DbName == "" {
		dbName := os.Getenv("EXPLORER_DB_NAME")
		env.DbName = dbName
	}

	if env.DbPassword == "" {
		dbPassword := os.Getenv("EXPLORER_DB_PASSWORD")
		env.DbPassword = dbPassword
	}

	if env.BaseCoin == "" {
		baseCoin := os.Getenv("EXPLORER_BASE_COIN")
		env.BaseCoin = baseCoin
	}

	return &env
}
