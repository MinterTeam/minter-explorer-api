package stake

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
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

// Get list of stakes by Minter address
func (repository Repository) GetByAddress(address string) []*models.Stake {
	var stakes []*models.Stake

	err := repository.db.Model(&stakes).
		Column("Coin", "OwnerAddress._").
		Where("owner_address.address = ?", address).
		Select()

	helpers.CheckErr(err)

	return stakes
}

// Get paginated list of stakes by Minter address
func (repository Repository) GetPaginatedByAddress(address string, pagination *tools.Pagination) ([]models.Stake, error) {
	var stakes []models.Stake
	var err error

	pagination.Total, err = repository.db.Model(&stakes).
		Column("Coin.symbol", "Validator.public_key", "OwnerAddress._").
		Column("Validator.name", "Validator.description", "Validator.icon_url", "Validator.site_url").
		Where("owner_address.address = ?", address).
		Apply(pagination.Filter).
		SelectAndCount()

	return stakes, err
}

// Get total delegated bip value
func (repository Repository) GetSumInBipValue() (string, error) {
	var sum string
	err := repository.db.Model(&models.Stake{}).ColumnExpr("SUM(bip_value)").Select(&sum)
	return sum, err
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
