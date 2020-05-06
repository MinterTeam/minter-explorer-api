package core

import (
	"log"
	"os"
	"strconv"
)

type Environment struct {
	DbName          string
	DbUser          string
	DbPassword      string
	DbPoolSize      int
	DbHost          string
	DbPort          string
	BaseCoin        string
	ServerPort      string
	IsDebug         bool
	WsServer        string
	WsBlocksChannel string
	BipdevApiHost   string
	NodeApiHost     string
}

func NewEnvironment() *Environment {
	dbPoolSize, err := strconv.ParseInt(os.Getenv("DB_POOL_SIZE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	env := Environment{
		DbName:          os.Getenv("DB_NAME"),
		DbUser:          os.Getenv("DB_USER"),
		DbPassword:      os.Getenv("DB_PASSWORD"),
		DbPort:          os.Getenv("DB_PORT"),
		DbPoolSize:      int(dbPoolSize),
		DbHost:          os.Getenv("DB_HOST"),
		BaseCoin:        os.Getenv("APP_BASE_COIN"),
		ServerPort:      os.Getenv("EXPLORER_PORT"),
		IsDebug:         os.Getenv("EXPLORER_DEBUG") == "1",
		WsServer:        os.Getenv("CENTRIFUGO_LINK"),
		WsBlocksChannel: os.Getenv("CENTRIFUGO_BLOCK_CHANNEL"),
		BipdevApiHost:   os.Getenv("BIP_DEV_API"),
		NodeApiHost:     os.Getenv("NODE_API"),
	}

	return &env
}
