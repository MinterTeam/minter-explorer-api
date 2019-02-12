package core

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/go-pg/pg"
)

type Explorer struct {
	CoinRepository  coins.Repository
	BlockRepository blocks.Repository
}

func NewExplorer(db *pg.DB) *Explorer {
	return &Explorer{
		CoinRepository:  *coins.NewRepository(db),
		BlockRepository: *blocks.NewRepository(db),
	}
}
