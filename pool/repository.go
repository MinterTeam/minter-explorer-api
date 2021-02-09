package pool

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
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
		Apply(filter.Filter("liquidity_pool__token", "liquidity_pool__first_coin", "liquidity_pool__second_coin")).
		Apply(pagination.Filter).
		SelectAndCount()

	return providers, err
}

func (r *Repository) FindRoutePath(filter SelectByCoinsFilter) ([]models.LiquidityPool, error) {
	fromCoinId, err := r.coinRepository.FindIdBySymbol(filter.Coin0)
	if err != nil {
		return nil, err
	}

	toCoinId, err := r.coinRepository.FindIdBySymbol(filter.Coin1)
	if err != nil {
		return nil, err
	}

	var path string
	_, err = r.db.QueryOne(&path, `WITH RECURSIVE search_graph(first_coin_id, second_coin_id, depth, path) AS (      
        SELECT g.first_coin_id, g.second_coin_id, 1 as depth, ARRAY[g.id] as path FROM liquidity_pools AS g WHERE first_coin_id = ?      
      	UNION ALL      
        SELECT g.first_coin_id, g.second_coin_id, sg.depth + 1 as depth, path || g.id as path
        FROM liquidity_pools AS g, search_graph AS sg      
        WHERE  g.first_coin_id = sg.second_coin_id AND (g.id <> ALL(sg.path)) AND sg.depth <= 3
	) SELECT path FROM search_graph where second_coin_id = ? order by depth limit 1;`, fromCoinId, toCoinId)

	if err != nil {
		return nil, err
	}

	var ids []uint64
	if err := json.Unmarshal([]byte(`[`+path[1:len(path)-1]+`]`), &ids); err != nil {
		return nil, err
	}

	var pools []models.LiquidityPool
	err = r.db.Model(&pools).
		Relation("Token").
		Relation("FirstCoin").
		Relation("SecondCoin").
		WhereIn("liquidity_pool.id IN (?)", ids).
		Select()

	return pools, err
}
