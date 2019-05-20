package coins

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/go-pg/pg"
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

	err := repository.DB.Model(&coins).Column("crr", "volume", "reserve_balance", "name", "symbol").
		Where("deleted_at IS NULL").Select()

	helpers.CheckErr(err)

	return coins
}

// Get coin detail by symbol
func (repository *Repository) GetBySymbol(symbol string) []models.Coin {
	var coins []models.Coin

	err := repository.DB.Model(&coins).Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Where("deleted_at IS NULL").
		Column("crr", "volume", "reserve_balance", "name", "symbol").Select()
	helpers.CheckErr(err)

	return coins
}

func (repository *Repository) GetCustomCoinsReserveSum() (string, error) {
	var sum string
	err := repository.DB.
		Model(&models.Coin{}).
		ColumnExpr("SUM(reserve_balance)").
		Where("symbol != ?", repository.baseCoinSymbol).
		Select(&sum)
	return sum, err
}
