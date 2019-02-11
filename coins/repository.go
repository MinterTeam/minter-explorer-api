package coins

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
)

type Repository struct {
	DB *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// Get list of coins
func (repository *Repository) GetCoins() []models.Coin {
	var coins []models.Coin

	// fetch data
	err := repository.DB.Model(&coins).Column("crr", "volume", "reserve_balance", "name", "symbol").Select()

	helpers.CheckErr(err)

	return coins
}

// Get coin detail by symbol
func (repository *Repository) GetBySymbol(symbol string) []models.Coin {
	var coins []models.Coin

	// fetch data
	err := repository.DB.Model(&coins).Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Column("crr", "volume", "reserve_balance", "name", "symbol").Select()

	helpers.CheckErr(err)

	return coins
}
