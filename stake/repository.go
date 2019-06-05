package stake

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/go-pg/pg"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get paginated list of stakes by Minter address
func (repository Repository) GetByAddress(address string, pagination *tools.Pagination) []models.Stake {
	var stakes []models.Stake
	var err error

	pagination.Total, err = repository.db.Model(&stakes).
		Column("Coin.symbol", "Validator.public_key", "OwnerAddress._").
		Where("owner_address.address = ?", address).
		Apply(pagination.Filter).
		SelectAndCount()

	helpers.CheckErr(err)

	return stakes
}

// Get total delegated sum by address
func (repository Repository) GetSumInBipValueByAddress(address string) (string, error) {
	var sum string
	err := repository.db.Model(&models.Stake{}).
		Column("OwnerAddress._").
		ColumnExpr("SUM(bip_value)").
		Where("owner_address.address = ?", address).
		Select(&sum)

	return sum, err
}
