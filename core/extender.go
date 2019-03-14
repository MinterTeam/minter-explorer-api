package core

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/centrifugal/centrifuge-go"
)

type Extender struct {
	client *centrifuge.Client
}

// create new extender connection
func NewExtender(explorer *Explorer) *Extender {
	c := centrifuge.New(explorer.Environment.WsServer, centrifuge.DefaultConfig())

	err := c.Connect()
	helpers.CheckErr(err)

	return &Extender{c}
}

// create new subscription to channel
func (e *Extender) CreateSubscription(channel string) *centrifuge.Subscription {
	sub, err := e.client.NewSubscription(channel)
	helpers.CheckErr(err)

	return sub
}

// subscribe to channel
func (e Extender) Subscribe(sub *centrifuge.Subscription) {
	err := sub.Subscribe()
	helpers.CheckErr(err)
}

// close ws connection
func (e *Extender) Close() {
	err := e.client.Close()
	helpers.CheckErr(err)
}
