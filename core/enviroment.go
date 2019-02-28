package core

import (
	"github.com/MinterTeam/minter-explorer-api/core/config"
)

type Environment struct {
	DbName     string
	DbUser     string
	DbPassword string
	BaseCoin   string
}

func NewEnvironment() *Environment {
	cfg := config.NewViperConfig()

	env := Environment{
		DbName:     cfg.GetString("database.name"),
		DbUser:     cfg.GetString("database.user"),
		DbPassword: cfg.GetString("database.password"),
		BaseCoin:   cfg.GetString("baseCoin"),
	}

	return &env
}
