package main

import (
	"github.com/MinterTeam/minter-explorer-api/v2/api"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/database"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/metrics"
	"github.com/joho/godotenv"
	"log"
)

func main() {
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
	go explorer.MarketService.Run()

	// create ws extender
	extender := core.NewExtenderWsClient(explorer)
	defer extender.Close()

	// subscribe to channel and add cache handler
	sub := extender.CreateSubscription(explorer.Environment.WsBlocksChannel)
	sub.OnPublish(explorer.Cache)
	extender.Subscribe(sub)

	// TODO: refactor
	// create ws extender for metrics
	extender2 := core.NewExtenderWsClient(explorer)
	defer extender2.Close()
	metricSub := extender2.CreateSubscription(explorer.Environment.WsBlocksChannel)
	metricSub.OnPublish(metrics.NewLastBlockMetric())
	extender2.Subscribe(metricSub)

	// run api
	api.Run(db, explorer)
}
