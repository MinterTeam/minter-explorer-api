package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"math/big"
	"time"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindByCoins(filter SelectByCoinsFilter) (models.LiquidityPool, error) {
	var pool models.LiquidityPool

	err := r.db.Model(&pool).
		Relation("Token").
		Relation("FirstCoin").
		Relation("SecondCoin").
		Apply(filter.Filter("token", "first_coin", "second_coin")).
		Order("liquidity_bip DESC").
		First()

	return pool, err
}

func (r *Repository) FindProvider(filter SelectByCoinsFilter, address string) (models.AddressLiquidityPool, error) {
	var provider models.AddressLiquidityPool

	err := r.db.Model(&provider).
		Relation("Address").
		Relation("LiquidityPool").
		Relation("LiquidityPool.Token").
		Relation("LiquidityPool.FirstCoin").
		Relation("LiquidityPool.SecondCoin").
		Where("address.address = ?", address).
		Apply(filter.Filter("liquidity_pool__token", "liquidity_pool__first_coin", "liquidity_pool__second_coin")).
		First()

	return provider, err
}

func (r *Repository) GetPools(filter SelectPoolsFilter, pagination *tools.Pagination) (pool []models.LiquidityPool, err error) {
	pagination.Total, err = r.db.Model(&pool).
		Relation("Token").
		Relation("FirstCoin").
		Relation("SecondCoin").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("liquidity_bip DESC").
		SelectAndCount()

	return pool, err
}

func (r *Repository) GetPoolsByProvider(provider string, pagination *tools.Pagination) (pools []models.AddressLiquidityPool, err error) {
	pagination.Total, err = r.db.Model(&pools).
		Relation("Address").
		Relation("LiquidityPool").
		Relation("LiquidityPool.Token").
		Relation("LiquidityPool.FirstCoin").
		Relation("LiquidityPool.SecondCoin").
		Where("address.address = ?", provider).
		Apply(pagination.Filter).
		OrderExpr(`(("address_liquidity_pool"."liquidity" / "liquidity_pool"."liquidity") * "liquidity_pool"."liquidity_bip" ) desc`).
		SelectAndCount()

	return pools, err
}

func (r *Repository) GetProviders(filter SelectByCoinsFilter, pagination *tools.Pagination) (providers []models.AddressLiquidityPool, err error) {
	pagination.Total, err = r.db.Model(&providers).
		Relation("Address").
		Relation("LiquidityPool").
		Relation("LiquidityPool.Token").
		Relation("LiquidityPool.FirstCoin").
		Relation("LiquidityPool.SecondCoin").
		OrderExpr(`(("address_liquidity_pool"."liquidity" / "liquidity_pool"."liquidity") * "liquidity_pool"."liquidity_bip" ) desc`).
		Apply(filter.Filter("liquidity_pool__token", "liquidity_pool__first_coin", "liquidity_pool__second_coin")).
		Apply(pagination.Filter).
		SelectAndCount()

	return providers, err
}

func (r *Repository) GetAll() (pools []models.LiquidityPool, err error) {
	err = r.db.Model(&pools).
		Relation("FirstCoin").
		Relation("SecondCoin").
		Relation("Token").
		Order("liquidity_bip DESC").
		Select()

	return pools, err
}

func (r *Repository) GetTracked() (pools []models.LiquidityPool, err error) {
	err = r.db.Model(&pools).
		Relation("FirstCoin").
		Relation("SecondCoin").
		Relation("Token").
		Where("first_coin_id in (select coin_id from token_contracts)").
		Where("second_coin_id in (select coin_id from token_contracts)").
		Order("liquidity_bip DESC").
		Select()

	return pools, err
}

func (r *Repository) Find(from, to uint64) (models.LiquidityPool, error) {
	var pool models.LiquidityPool

	err := r.db.Model(&pool).
		WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.Where("first_coin_id = ?", from).Where("second_coin_id = ?", to), nil
		}).
		WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
			return query.Where("first_coin_id = ?", to).Where("second_coin_id = ?", from), nil
		}).
		Order("id").
		First()

	return pool, err
}

func (r *Repository) GetPoolsCoins() (coins []models.Coin, err error) {
	err = r.db.Model(&coins).
		Where(`exists (select * from liquidity_pools where first_coin_id = "coin"."id" or second_coin_id = "coin"."id")`).
		Order("reserve DESC").
		Select()

	return coins, err
}

func (r *Repository) GetPoolTradesVolume(pool models.LiquidityPool, scale string, startTime *time.Time) (trades []tradeVolume, err error) {
	q := r.db.Model(&models.LiquidityPoolTrade{}).
		ColumnExpr("sum(first_coin_volume) as first_coin_volume").
		ColumnExpr("sum(second_coin_volume) as second_coin_volume").
		ColumnExpr("date_trunc(?, created_at) as date", scale).
		Where("liquidity_pool_id = ?", pool.Id).
		Group("date").
		Order("date DESC")

	if startTime != nil {
		q.Where("created_at > ?", startTime.Format(time.RFC3339))
	}

	err = q.Select(&trades)

	return trades, err
}

func (r *Repository) GetPoolTradeVolumeByTimeRange(pool models.LiquidityPool, startTime time.Time) (*tradeVolume, error) {
	tv := new(tradeVolume)
	count, err := r.db.Model(&models.LiquidityPoolTrade{}).
		ColumnExpr("liquidity_pool_id as pool_id").
		ColumnExpr("sum(first_coin_volume) as first_coin_volume").
		ColumnExpr("sum(second_coin_volume) as second_coin_volume").
		Group("liquidity_pool_id").
		Where("liquidity_pool_id = ?", pool.Id).
		Where("created_at > ?", startTime.Format(time.RFC3339)).
		SelectAndCount(tv)

	if count == 0 {
		return nil, nil
	}

	return tv, err
}

func (r *Repository) GetPoolsTradeVolumeByTimeRange(pools []models.LiquidityPool, startTime time.Time) (tvs []tradeVolume, err error) {
	ids := make([]uint64, len(pools))
	for i, p := range pools {
		ids[i] = p.Id
	}

	count, err := r.db.Model(&models.LiquidityPoolTrade{}).
		ColumnExpr("liquidity_pool_id as pool_id").
		ColumnExpr("sum(first_coin_volume) as first_coin_volume").
		ColumnExpr("sum(second_coin_volume) as second_coin_volume").
		Group("liquidity_pool_id").
		Where("liquidity_pool_id in (?)", pg.In(ids)).
		Where("created_at > ?", startTime.Format(time.RFC3339)).
		SelectAndCount(&tvs)

	if count == 0 {
		return nil, nil
	}

	return tvs, err
}

func (r *Repository) GetPoolsCount() (count int, err error) {
	return r.db.Model(&models.LiquidityPool{}).Count()
}

func (r *Repository) GetTokenContractByCoinId(coinId uint64) (*models.TokenContract, error) {
	var tc models.TokenContract
	if err := r.db.Model(&tc).Where("coin_id = ? ", coinId).First(); err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *Repository) GetCoinsTradingVolume(scale string) ([]CoinTradingVolume, error) {
	var coinTradingVolumes []CoinTradingVolume

	tradingVolumes := r.db.Model(new(models.LiquidityPoolTrade)).
		ColumnExpr(`sum("liquidity_pool_trade"."first_coin_volume") as volume`).
		ColumnExpr(`liquidity_pools.first_coin_id as coin_id`).
		Join(`JOIN liquidity_pools on liquidity_pools.id = "liquidity_pool_trade"."liquidity_pool_id"`).
		Where("created_at > 'now'::timestamp - ?::interval", scale).
		Group("liquidity_pools.first_coin_id").
		UnionAll(
			r.db.Model(new(models.LiquidityPoolTrade)).
				ColumnExpr(`sum("liquidity_pool_trade"."second_coin_volume") as volume`).
				ColumnExpr(`liquidity_pools.second_coin_id as coin_id`).
				Join(`JOIN liquidity_pools on liquidity_pools.id = "liquidity_pool_trade"."liquidity_pool_id"`).
				Where("created_at > 'now'::timestamp - ?::interval", scale).
				Group("liquidity_pools.second_coin_id"),
		)

	err := r.db.Model().
		With("data", tradingVolumes).
		Table("data").
		ColumnExpr("sum(data.volume) as volume").
		ColumnExpr("data.coin_id").
		GroupExpr("data.coin_id").
		Select(&coinTradingVolumes)

	return coinTradingVolumes, err
}

func (r *Repository) GetTotalValueLocked() (*big.Int, error) {
	var tvlStr string
	err := r.db.Model(new(models.LiquidityPool)).ColumnExpr("sum(liquidity_bip)").Select(&tvlStr)
	if err != nil {
		return nil, err
	}

	return helpers.StringToBigInt(tvlStr), err
}
