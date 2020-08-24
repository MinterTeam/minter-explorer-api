package coins

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
)

type Repository struct {
	DB             *pg.DB
	baseCoinSymbol string
}

func NewRepository(db *pg.DB, baseCoinSymbol string) *Repository {
	return &Repository{
		DB:             db,
		baseCoinSymbol: baseCoinSymbol,
	}
}

// Get list of coins
func (repository *Repository) GetCoins() []models.Coin {
	var coins []models.Coin

	err := repository.DB.Model(&coins).
		Column("OwnerAddress").
		Where("deleted_at IS NULL").
		Order("reserve DESC").
		Select()

	helpers.CheckErr(err)

	return coins
}

// Get coin detail by symbol
func (repository *Repository) GetBySymbol(symbol string) []models.Coin {
	var coins []models.Coin

	err := repository.DB.Model(&coins).
		Column("Address").
		Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Where("deleted_at IS NULL").
		Order("reserve_balance DESC").
		Select()
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
		Where("symbol != ?", repository.baseCoinSymbol).
		Select(&data)

	return data, err
}
