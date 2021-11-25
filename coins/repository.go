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
	DB        *pg.DB
	baseModel *models.Coin
	coins     *sync.Map

	blockListCoinIds []uint64
}

var GlobalRepository *Repository

func NewRepository(db *pg.DB) *Repository {
	GlobalRepository = &Repository{
		DB:    db,
		coins: new(sync.Map),
	}

	return GlobalRepository
}

func (repository Repository) GetAll() []models.Coin {
	var coins []models.Coin
	err := repository.DB.Model(&coins).Relation("OwnerAddress").Where("deleted_at is null").Select()
	helpers.CheckErr(err)
	return coins
}

// Get list of coins
func (repository *Repository) GetCoins() []models.Coin {
	var coins []models.Coin

	allCoins := repository.DB.Model(new(models.Coin)).
		ColumnExpr("id").
		ColumnExpr("reserve").
		Where("deleted_at is null").
		UnionAll(repository.DB.Model(new(models.LiquidityPool)).
			ColumnExpr("first_coin_id as id").
			ColumnExpr("liquidity_bip as reserve")).
		UnionAll(repository.DB.Model(new(models.LiquidityPool)).
			ColumnExpr("second_coin_id as id").
			ColumnExpr("liquidity_bip as reserve"))

	err := repository.DB.Model().
		With("all_coins", allCoins).
		Table("all_coins").
		ColumnExpr("coins.*").
		ColumnExpr(`addresses.address as owner_address__address`).
		Join("JOIN coins ON coins.id = all_coins.id").
		Join("left outer JOIN addresses on addresses.id = coins.owner_address_id").
		Where("coins.id not in (?)", pg.In(repository.blockListCoinIds)).
		GroupExpr("coins.id, addresses.address").
		OrderExpr(`case when coins.id = 0 then 0 else 1 end`).
		OrderExpr("max(all_coins.reserve) desc").
		Select(&coins)

	helpers.CheckErr(err)

	return coins
}

// Get coin detail like symbol
func (repository *Repository) GetLikeSymbolAndVersion(symbol string, version *uint64) []models.Coin {
	var coins []models.Coin

	query := repository.DB.Model(&coins).
		Relation("OwnerAddress").
		Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Where("deleted_at IS NULL").
		Where(`"coin"."id" not in (?)`, pg.In(repository.blockListCoinIds)).
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
func (repository *Repository) GetBySymbolAndVersion(symbol string, version *uint64) []models.Coin {
	var coins []models.Coin

	query := repository.DB.Model(&coins).
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
func (repository *Repository) GetCustomCoinsStatusData() (CustomCoinsStatusData, error) {
	var data CustomCoinsStatusData

	err := repository.DB.
		Model(&models.Coin{}).
		ColumnExpr("SUM(reserve) as reserve_sum, COUNT(*) as count").
		Where("id != ?", 0).
		Select(&data)

	return data, err
}

func (repository *Repository) FindByID(id uint) (models.Coin, error) {
	var coin models.Coin

	if id == 0 && repository.baseModel != nil {
		return *repository.baseModel, nil
	}

	if c, ok := repository.coins.Load(id); ok {
		return c.(models.Coin), nil
	}

	err := repository.DB.Model(&coin).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Select()

	if id == 0 && repository.baseModel == nil {
		repository.baseModel = &coin
	}

	if err == nil {
		repository.coins.Store(id, coin)
	}

	return coin, err
}

func (repository Repository) FindByIdWithOwner(id uint) (models.Coin, error) {
	var coin models.Coin

	if id == 0 && repository.baseModel != nil {
		return *repository.baseModel, nil
	}

	err := repository.DB.Model(&coin).
		Relation("OwnerAddress").
		Where(`"coin"."id" = ?`, id).
		Where("deleted_at IS NULL").
		Select()

	if id == 0 && repository.baseModel == nil {
		repository.baseModel = &coin
	}

	return coin, err
}

func (repository *Repository) FindIdBySymbol(symbol string) (uint64, error) {
	if id, err := strconv.ParseUint(symbol, 10, 64); err != nil {
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(symbol)
		coins := repository.GetBySymbolAndVersion(symbol, &version)
		if len(coins) == 0 {
			return 0, pg.ErrNoRows
		}

		return uint64(coins[0].ID), nil
	} else {
		return id, nil
	}
}

func (repository *Repository) GetBySymbols(symbols []string) (coins []models.Coin, err error) {
	err = repository.DB.Model(&coins).Where("symbol in (?)", pg.In(symbols)).Select()
	return
}

func (repository *Repository) OnNewBlock(block blocks.Resource) {
	repository.fillCoinsMap()
}

func (repository *Repository) fillCoinsMap() {
	wg := &sync.WaitGroup{}
	for _, coin := range repository.GetAll() {
		wg.Add(1)
		go func(wg *sync.WaitGroup, coin models.Coin) {
			defer wg.Done()
			repository.coins.Store(uint64(coin.ID), coin)
		}(wg, coin)
	}
	wg.Wait()
}

func (repository *Repository) SetBlocklistCoinIds(ids []uint64) {
	repository.blockListCoinIds = ids
}
