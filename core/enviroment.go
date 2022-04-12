package core

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
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
	Basecoin        string
	ServerPort      string
	IsDebug         bool
	WsServer        string
	WsBlocksChannel string
	MarketHost      string
}

func NewEnvironment() *Environment {
	dbPoolSize, err := strconv.ParseInt(os.Getenv("DB_POOL_SIZE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	config.BaseCoinSymbol = os.Getenv("APP_BASE_COIN")
	config.SwapRouterProxyUrl = os.Getenv("SWAP_ROUTER_PROXY_URL")

	env := Environment{
		DbName:          os.Getenv("DB_NAME"),
		DbUser:          os.Getenv("DB_USER"),
		DbPassword:      os.Getenv("DB_PASSWORD"),
		DbPort:          os.Getenv("DB_PORT"),
		DbPoolSize:      int(dbPoolSize),
		DbHost:          os.Getenv("DB_HOST"),
		Basecoin:        os.Getenv("APP_BASE_COIN"),
		ServerPort:      os.Getenv("EXPLORER_PORT"),
		IsDebug:         os.Getenv("EXPLORER_DEBUG") == "1",
		WsServer:        os.Getenv("CENTRIFUGO_LINK"),
		WsBlocksChannel: os.Getenv("CENTRIFUGO_BLOCK_CHANNEL"),
		MarketHost:      os.Getenv("MARKET_HOST"),
	}

	return &env
}
