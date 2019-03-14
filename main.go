package main

import (
	"github.com/MinterTeam/minter-explorer-api/api"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/database"
)

func main() {
	// init environment
	env := core.NewEnvironment()

	// connect to database
	db := database.Connect(env)
	defer database.Close(db)

	// create explorer
	explorer := core.NewExplorer(db, env)

	// create ws extender
	extender := core.NewExtender(explorer)
	defer extender.Close()

	// subscribe to channel and add cache handler
	sub := extender.CreateSubscription(explorer.Environment.WsBlocksChannel)
	sub.OnPublish(explorer.Cache)
	extender.Subscribe(sub)

	// run api
	api.Run(db, explorer)
}
