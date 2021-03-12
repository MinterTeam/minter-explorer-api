package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Repository struct {
	db             *pg.DB
	coinRepository *coins.Repository
}

func NewRepository(db *pg.DB, coinRepository *coins.Repository) *Repository {
	return &Repository{db, coinRepository}
}

func (r *Repository) FindByCoins(filter SelectByCoinsFilter) (models.LiquidityPool, error) {
	var pool models.LiquidityPool

	err := r.db.Model(&pool).
		Relation("Token").
		Relation("FirstCoin").
		Relation("SecondCoin").
		Apply(filter.Filter("token", "first_coin", "second_coin")).
		Order("id").
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
		Order("id").
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
		Order(`address_liquidity_pool.liquidity DESC`).
		Apply(pagination.Filter).
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
		Order(`address_liquidity_pool.liquidity DESC`).
		Apply(filter.Filter("liquidity_pool__token", "liquidity_pool__first_coin", "liquidity_pool__second_coin")).
		Apply(pagination.Filter).
		SelectAndCount()

	return providers, err
}

func (r *Repository) GetAll() (pools []models.LiquidityPool, err error) {
	err = r.db.Model(&pools).
		Relation("FirstCoin").
		Relation("SecondCoin").
		Order("id").
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
