package coins

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
	"strconv"
	"sync"
)

type Repository struct {
	db        *pg.DB
	baseModel *models.Coin
	coins     *sync.Map

	blockListCoinIds []uint64
}

var GlobalRepository *Repository

func NewRepository(db *pg.DB) *Repository {
	GlobalRepository = &Repository{
		db:    db,
		coins: new(sync.Map),
	}

	return GlobalRepository
}

func (r Repository) GetAll() []models.Coin {
	var coins []models.Coin
	err := r.db.Model(&coins).Relation("OwnerAddress").Where("deleted_at is null").Select()
	helpers.CheckErr(err)
	return coins
}

// Get list of coins
func (r *Repository) GetCoins() []models.Coin {
	var coins []models.Coin

	allCoins := r.db.Model(new(models.Coin)).
		ColumnExpr("id").
		ColumnExpr("reserve").
		Where("deleted_at is null").
		UnionAll(r.db.Model(new(models.LiquidityPool)).
			ColumnExpr("first_coin_id as id").
			ColumnExpr("liquidity_bip as reserve")).
		UnionAll(r.db.Model(new(models.LiquidityPool)).
			ColumnExpr("second_coin_id as id").
			ColumnExpr("liquidity_bip as reserve"))

	err := r.db.Model().
		With("all_coins", allCoins).
		Table("all_coins").
		ColumnExpr("coins.*").
		ColumnExpr(`addresses.address as owner_address__address`).
		Join("JOIN coins ON coins.id = all_coins.id").
		Join("left outer JOIN addresses on addresses.id = coins.owner_address_id").
		Where("coins.id not in (?)", pg.In(r.blockListCoinIds)).
		GroupExpr("coins.id, addresses.address").
		OrderExpr(`case when coins.id = 0 then 0 else 1 end`).
		OrderExpr("max(all_coins.reserve) desc").
		Select(&coins)

	helpers.CheckErr(err)

	return coins
}

// Get coin detail like symbol
func (r *Repository) GetLikeSymbolAndVersion(symbol string, version *uint64) []models.Coin {
	var coins []models.Coin

	query := r.db.Model(&coins).
		Relation("OwnerAddress").
		Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Where("deleted_at IS NULL").
		Where(`"coin"."id" not in (?)`, pg.In(r.blockListCoinIds)).
		OrderExpr(`case when "coin"."id" = 0 then 0 else 1 end`).
		Order("reserve DESC")

	if version != nil {
		query.Where("version = ?", version)
	}

	err := query.Select()
	helpers.CheckErr(err)

	return coins
}

// Get coin detail by symbol
func (r *Repository) GetBySymbolAndVersion(symbol string, version *uint64) []models.Coin {
	var coins []models.Coin

	query := r.db.Model(&coins).
		Relation("OwnerAddress").
		Where("symbol = ?", symbol).
		Where("deleted_at IS NULL").
		OrderExpr(`case when "coin"."id" = 0 then 0 else 1 end`).
		Order("reserve DESC")

	if version != nil {
		query.Where("version = ?", version)
	}

	err := query.Select()
	helpers.CheckErr(err)

	return coins
}

type CustomCoinsStatusData struct {
	ReserveSum string
	Count      uint
}

// Get custom coins data for status page
func (r *Repository) GetCustomCoinsStatusData() (CustomCoinsStatusData, error) {
	var data CustomCoinsStatusData

	err := r.db.
		Model(&models.Coin{}).
		ColumnExpr("SUM(reserve) as reserve_sum, COUNT(*) as count").
		Where("id != ?", 0).
		Select(&data)

	return data, err
}

func (r *Repository) FindByID(id uint) (models.Coin, error) {
	var coin models.Coin

	//if id == 0 && r.baseModel != nil {
	//	return *r.baseModel, nil
	//}

	if c, ok := r.coins.Load(id); ok {
		return c.(models.Coin), nil
	}

	err := r.db.Model(&coin).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Select()

	//if id == 0 && r.baseModel == nil {
	//	r.baseModel = &coin
	//}

	if err == nil {
		r.coins.Store(id, coin)
	}

	return coin, err
}

func (r Repository) FindByIdWithOwner(id uint) (models.Coin, error) {
	var coin models.Coin

	if id == 0 && r.baseModel != nil {
		return *r.baseModel, nil
	}

	err := r.db.Model(&coin).
		Relation("OwnerAddress").
		Where(`"coin"."id" = ?`, id).
		Where("deleted_at IS NULL").
		Select()

	if id == 0 && r.baseModel == nil {
		r.baseModel = &coin
	}

	return coin, err
}

func (r *Repository) FindIdBySymbol(symbol string) (uint64, error) {
	if id, err := strconv.ParseUint(symbol, 10, 64); err != nil {
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(symbol)
		coins := r.GetBySymbolAndVersion(symbol, &version)
		if len(coins) == 0 {
			return 0, pg.ErrNoRows
		}

		return uint64(coins[0].ID), nil
	} else {
		return id, nil
	}
}

func (r *Repository) GetBySymbols(symbols []string) (coins []models.Coin, err error) {
	err = r.db.Model(&coins).Where("symbol in (?)", pg.In(symbols)).Select()
	return
}

func (r *Repository) OnNewBlock(block blocks.Resource) {
	r.fillCoinsMap()
}

func (r *Repository) fillCoinsMap() {
	wg := &sync.WaitGroup{}
	for _, coin := range r.GetAll() {
		wg.Add(1)
		go func(wg *sync.WaitGroup, coin models.Coin) {
			defer wg.Done()
			r.coins.Store(uint64(coin.ID), coin)
		}(wg, coin)
	}
	wg.Wait()
}

func (r *Repository) SetBlocklistCoinIds(ids []uint64) {
	r.blockListCoinIds = ids
}

func (r *Repository) GetVerifiedCoins() (coins []models.Coin, err error) {
	err = r.db.Model(&coins).Where("id in (select coin_id from token_contracts)").Select()
	return
}