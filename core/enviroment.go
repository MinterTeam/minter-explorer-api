package core

import (
	"flag"
	"os"
)

type Environment struct {
	DbName     string
	DbUser     string
	DbPassword string
}

func NewEnvironment() *Environment {
	dbName := flag.String("db_name", "", "DB name")
	dbUser := flag.String("db_user", "", "DB user")
	dbPassword := flag.String("db_password", "", "DB password")
	flag.Parse()

	env := Environment{
		DbName:     *dbName,
		DbUser:     *dbUser,
		DbPassword: *dbPassword,
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

	return &env
}
