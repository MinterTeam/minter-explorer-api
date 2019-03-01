package core

import (
	"github.com/MinterTeam/minter-explorer-api/core/config"
)

type Environment struct {
	DbName     string
	DbUser     string
	DbPassword string
	DbPoolSize int
	BaseCoin   string
	ServerPort string
}

func NewEnvironment() *Environment {
	cfg := config.NewViperConfig()

	env := Environment{
		DbName:     cfg.GetString("database.name"),
		DbUser:     cfg.GetString("database.user"),
		DbPassword: cfg.GetString("database.password"),
		DbPoolSize: cfg.GetInt("database.poolSize"),
		BaseCoin:   cfg.GetString("baseCoin"),
		ServerPort: cfg.GetString("server.port"),
	}

	return &env
}
