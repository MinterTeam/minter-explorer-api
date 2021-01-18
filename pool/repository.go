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

func (r *Repository) FindByCoins(coin0 uint64, coin1 uint64) (models.LiquidityPool, error) {
	var pool models.LiquidityPool

	err := r.db.Model(&pool).
		Relation("FirstCoin").
		Relation("SecondCoin").
		Where("first_coin_id = ?", coin0).
		Where("second_coin_id = ?", coin1).
		First()

	return pool, err
}

func (r *Repository) FindProvider(coin0 uint64, coin1 uint64, address string) (models.AddressLiquidityPool, error) {
	var provider models.AddressLiquidityPool

	err := r.db.Model(&provider).
		Relation("Address").
		Relation("LiquidityPool").
		Where("address.address = ?", address).
		Where("liquidity_pool.first_coin_id = ?", coin0).
		Where("liquidity_pool.second_coin_id = ?", coin1).
		First()

	return provider, err
}

func (r *Repository) GetPools(filter SelectPoolsFilter, pagination *tools.Pagination) (provider []models.LiquidityPool, err error) {
	pagination.Total, err = r.db.Model(&provider).
		Relation("FirstCoin").
		Relation("SecondCoin").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		SelectAndCount()

	return provider, err
}
