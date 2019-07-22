package core

import (
	"github.com/MinterTeam/minter-explorer-api/address"
	"github.com/MinterTeam/minter-explorer-api/bipdev"
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/invalid_transaction"
	"github.com/MinterTeam/minter-explorer-api/reward"
	"github.com/MinterTeam/minter-explorer-api/slash"
	"github.com/MinterTeam/minter-explorer-api/stake"
	"github.com/MinterTeam/minter-explorer-api/tools/cache"
	"github.com/MinterTeam/minter-explorer-api/tools/market"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/go-pg/pg"
)

type Explorer struct {
	CoinRepository               coins.Repository
	BlockRepository              blocks.Repository
	AddressRepository            address.Repository
	TransactionRepository        transaction.Repository
	InvalidTransactionRepository invalid_transaction.Repository
	RewardRepository             reward.Repository
	SlashRepository              slash.Repository
	ValidatorRepository          validator.Repository
	StakeRepository              stake.Repository
	Environment                  Environment
	Cache                        *cache.ExplorerCache
	MarketService                *market.Service
}

func NewExplorer(db *pg.DB, env *Environment) *Explorer {
	return &Explorer{
		CoinRepository:               *coins.NewRepository(db, env.BaseCoin),
		BlockRepository:              *blocks.NewRepository(db),
		AddressRepository:            *address.NewRepository(db),
		TransactionRepository:        *transaction.NewRepository(db),
		InvalidTransactionRepository: *invalid_transaction.NewRepository(db),
		RewardRepository:             *reward.NewRepository(db),
		SlashRepository:              *slash.NewRepository(db),
		ValidatorRepository:          *validator.NewRepository(db),
		StakeRepository:              *stake.NewRepository(db),
		Environment:                  *env,
		Cache:                        cache.NewCache(),
		MarketService:                market.NewService(bipdev.NewApi(env.BipdevApiHost), env.BaseCoin),
	}
}
