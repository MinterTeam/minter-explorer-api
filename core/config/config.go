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

const (
	MUSD = uint64(2024)
)

var (
	StableCoinIds = []uint64{1993, 1994, 1995, 1996, 1997, 1998, 1999, 2000, 2024}
)
