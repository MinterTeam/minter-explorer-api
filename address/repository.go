package address

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get address model by address
func (r Repository) GetByAddress(minterAddress string) *models.Address {
	var address models.Address

	err := r.db.Model(&address).
		Relation("Balances").
		Relation("Balances.Coin").
		Where("address = ?", minterAddress).
		Select()

	if err != nil {
		return nil
	}

	return &address
}

// Get list of addresses models
func (r Repository) GetByAddresses(minterAddresses []string) []*models.Address {
	var addresses []*models.Address

	err := r.db.Model(&addresses).
		Relation("Balances").
		Relation("Balances.Coin").
		WhereIn("address IN (?)", minterAddresses).
		Select()

	helpers.CheckErr(err)

	return addresses
}

func (r Repository) GetNonZeroAddressesCount() (count uint64, err error) {
	err = r.db.Model(new(models.Balance)).
		ColumnExpr("count (DISTINCT address_id)").
		Select(&count)

	return count, err
}
