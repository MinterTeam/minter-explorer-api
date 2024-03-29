package main

import (
	"github.com/MinterTeam/minter-explorer-api/v2/api"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/ws"
	"github.com/MinterTeam/minter-explorer-api/v2/database"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/metrics"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/recovery"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	// Init Logger
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found")
	}

	// init environment
	env := core.NewEnvironment()

	// connect to database
	db := database.Connect(env)
	defer database.Close(db)

	// create explorer
	explorer := core.NewExplorer(db, env)

	// run market price update
	go recovery.SafeGo(explorer.MarketService.Run)
	go recovery.SafeGo(explorer.PoolService.RunPoolUpdater)
	go recovery.SafeGo(explorer.PoolService.RunTradingVolumeUpdater)

	// create ws extender
	extender := ws.NewExtenderWsClient(explorer)
	defer extender.Close()

	// create ws channel handler
	blocksChannelHandler := ws.NewBlocksChannelHandler()
	blocksChannelHandler.AddSubscriber(explorer.Cache)
	blocksChannelHandler.AddSubscriber(metrics.NewLastBlockMetric())
	blocksChannelHandler.AddSubscriber(explorer.PoolService)
	blocksChannelHandler.AddSubscriber(explorer.CoinRepository)

	// subscribe to channel and add cache handler
	sub := extender.CreateSubscription(explorer.Environment.WsBlocksChannel)
	sub.OnPublish(blocksChannelHandler)
	extender.Subscribe(sub)

	api.Run(db, explorer)
}
