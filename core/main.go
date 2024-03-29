package core

import (
	"github.com/MinterTeam/minter-explorer-api/v2/address"
	"github.com/MinterTeam/minter-explorer-api/v2/balance"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/check"
	"github.com/MinterTeam/minter-explorer-api/v2/coingecko"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/invalid_transaction"
	"github.com/MinterTeam/minter-explorer-api/v2/order"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-api/v2/reward"
	"github.com/MinterTeam/minter-explorer-api/v2/services"
	"github.com/MinterTeam/minter-explorer-api/v2/slash"
	"github.com/MinterTeam/minter-explorer-api/v2/stake"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/cache"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/market"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-api/v2/unbond"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/go-pg/pg/v10"
)

type Explorer struct {
	CoinRepository               *coins.Repository
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
	BalanceService               *balance.Service
	ValidatorService             *services.ValidatorService
	TransactionService           *transaction.Service
	UnbondRepository             *unbond.Repository
	StakeService                 *stake.Service
	CheckRepository              *check.Repository
	PoolRepository               *pool.Repository
	PoolService                  *pool.Service
	SwapService                  *services.SwapService
	AddressService               *address.Service
	CoinService                  *coins.Service
	OrderRepository              *order.Repository
}

func NewExplorer(db *pg.DB, env *Environment) *Explorer {
	marketService := market.NewService(coingecko.NewService(env.MarketHost), env.Basecoin)
	blockRepository := *blocks.NewRepository(db)
	validatorRepository := validator.NewRepository(db)
	stakeRepository := stake.NewRepository(db)
	coinRepository := coins.NewRepository(db)
	poolRepository := pool.NewRepository(db)
	addressRepository := address.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)
	coinService := coins.NewService(coinRepository)
	poolService := pool.NewService(poolRepository, coinService)
	services.Swap = services.NewSwapService(poolService)
	cacheService := cache.NewCache(blockRepository.GetLastBlock())
	transactionService := transaction.NewService(coinRepository, cacheService, transactionRepository)
	balanceService := balance.NewService(env.Basecoin, poolService, services.Swap)
	addressService := address.NewService(addressRepository, stakeRepository, balanceService)

	return &Explorer{
		BlockRepository:              blockRepository,
		CoinRepository:               coinRepository,
		AddressRepository:            *addressRepository,
		TransactionRepository:        *transactionRepository,
		InvalidTransactionRepository: *invalid_transaction.NewRepository(db),
		RewardRepository:             *reward.NewRepository(db),
		SlashRepository:              *slash.NewRepository(db),
		ValidatorRepository:          *validatorRepository,
		StakeRepository:              *stakeRepository,
		Environment:                  *env,
		Cache:                        cacheService,
		MarketService:                marketService,
		TransactionService:           transactionService,
		BalanceService:               balanceService,
		ValidatorService:             services.NewValidatorService(stakeRepository, cacheService),
		UnbondRepository:             unbond.NewRepository(db),
		StakeService:                 stake.NewService(stakeRepository),
		CheckRepository:              check.NewRepository(db),
		PoolRepository:               poolRepository,
		PoolService:                  poolService,
		SwapService:                  services.Swap,
		AddressService:               addressService,
		CoinService:                  coinService,
		OrderRepository:              order.NewRepository(db),
	}
}
