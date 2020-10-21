package ws

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/centrifugal/centrifuge-go"
)

type BlocksChannelHandler struct {
	Subscribers []NewBlockSubscriber
}

type NewBlockSubscriber interface {
	OnNewBlock(blocks.Resource)
}

// Constructor for blocks event channel handler
func NewBlocksChannelHandler() *BlocksChannelHandler {
	return &BlocksChannelHandler{
		Subscribers: make([]NewBlockSubscriber, 0),
	}
}

// Add new subscriber for channel
func (b *BlocksChannelHandler) AddSubscriber(sub NewBlockSubscriber) {
	b.Subscribers = append(b.Subscribers, sub)
}

// Handle new block from ws and publish him to all subscribers
func (b *BlocksChannelHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	var block blocks.Resource
	err := json.Unmarshal(e.Data, &block)
	helpers.CheckErr(err)

	for _, sub := range b.Subscribers {
		go sub.OnNewBlock(block)
	}
}
