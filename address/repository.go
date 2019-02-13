package address

import (
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
func (repository Repository) GetByAddresses(minterAddresses []string) *[]models.Address {
	var addresses []models.Address

	addressesList := make([]interface{}, len(minterAddresses))
	for i, address := range minterAddresses {
		addressesList[i] = address
	}

	err := repository.DB.Model(&addresses).Column("Balances", "Balances.Coin").
		WhereIn("address IN (?)", addressesList...).Select()
	helpers.CheckErr(err)

	return &addresses
}