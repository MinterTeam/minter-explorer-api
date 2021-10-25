package config

const (
	DefaultStatisticsScale       = "day"
	DefaultStatisticsDayDelta    = -14
	NetworkActivePeriod          = 90
	SlowBlocksMaxTimeInSec       = 6
	DefaultPaginationLimit       = 50
	MaxPaginationOffset          = 5000000
	MaxPaginationLimit           = 150
	MarketPriceUpdatePeriodInMin = 2
	MaxDelegatorCount            = 1000
)

var (
	BaseCoinSymbol string
	PriceCoinId    uint64 = 1993 // TODO: move to env
)
