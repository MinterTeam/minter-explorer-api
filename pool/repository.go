package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
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
		Relation("FirstCoin").
		Relation("SecondCoin").
		Apply(filter.Filter("first_coin", "second_coin")).
		First()

	return pool, err
}

func (r *Repository) FindProvider(filter SelectByCoinsFilter, address string) (models.AddressLiquidityPool, error) {
	var provider models.AddressLiquidityPool

	err := r.db.Model(&provider).
		Relation("Address").
		Relation("LiquidityPool").
		Relation("LiquidityPool.FirstCoin").
		Relation("LiquidityPool.SecondCoin").
		Where("address.address = ?", address).
		Apply(filter.Filter("liquidity_pool__first_coin", "liquidity_pool__second_coin")).
		First()

	return provider, err
}

func (r *Repository) GetPools(filter SelectPoolsFilter, pagination *tools.Pagination) (pool []models.LiquidityPool, err error) {
	pagination.Total, err = r.db.Model(&pool).
		Relation("FirstCoin").
		Relation("SecondCoin").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		SelectAndCount()

	return pool, err
}

func (r *Repository) GetProviders(filter SelectByCoinsFilter, pagination *tools.Pagination) (providers []models.AddressLiquidityPool, err error) {
	pagination.Total, err = r.db.Model(&providers).
		Relation("Address").
		Relation("LiquidityPool").
		Relation("LiquidityPool.FirstCoin").
		Relation("LiquidityPool.SecondCoin").
		Apply(filter.Filter("liquidity_pool__first_coin", "liquidity_pool__second_coin")).
		Apply(pagination.Filter).
		SelectAndCount()

	return providers, err
}
