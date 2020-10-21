package ws

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/centrifugal/centrifuge-go"
)

type ExtenderWsClient struct {
	client *centrifuge.Client
}

// create new extender connection
func NewExtenderWsClient(explorer *core.Explorer) *ExtenderWsClient {
	c := centrifuge.New(explorer.Environment.WsServer, centrifuge.DefaultConfig())

	err := c.Connect()
	helpers.CheckErr(err)

	return &ExtenderWsClient{c}
}

// create new subscription to channel
func (e *ExtenderWsClient) CreateSubscription(channel string) *centrifuge.Subscription {
	sub, err := e.client.NewSubscription(channel)
	helpers.CheckErr(err)

	return sub
}

// subscribe to channel
func (e ExtenderWsClient) Subscribe(sub *centrifuge.Subscription) {
	err := sub.Subscribe()
	helpers.CheckErr(err)
}

// close ws connection
func (e *ExtenderWsClient) Close() {
	err := e.client.Close()
	helpers.CheckErr(err)
}
