package address

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
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

// Get address model by address
func (repository Repository) GetByAddress(minterAddress string) *models.Address {
	var address models.Address

	err := repository.DB.Model(&address).Column("Balances", "Balances.Coin").
		Where("address = ?", minterAddress).Select()
	if err != nil {
		return nil
	}

	return &address
}

// Get list of addresses models
func (repository Repository) GetByAddresses(minterAddresses []string) []models.Address {
	var addresses []models.Address

	err := repository.DB.Model(&addresses).Column("Balances", "Balances.Coin").
		WhereIn("address IN (?)", pg.In(minterAddresses)).Select()

	helpers.CheckErr(err)

	return addresses
}
