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
		Apply(filter.Filter).
		First()

	return pool, err
}

func (r *Repository) FindProvider(filter SelectByCoinsFilter, address string) (models.AddressLiquidityPool, error) {
	var provider models.AddressLiquidityPool

	err := r.db.Model(&provider).
		Relation("Address").
		Relation("LiquidityPool").
		Join("JOIN coins as first_coin").
		JoinOn("first_coin.id = liquidity_pool.first_coin_id").
		Join("JOIN coins as second_coin").
		JoinOn("second_coin.id = liquidity_pool.second_coin_id").
		Where("address.address = ?", address).
		Apply(filter.Filter).
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
