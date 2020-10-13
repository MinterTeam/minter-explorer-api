package stake

import (
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
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
func (repository Repository) GetAllByAddress(address string) ([]models.Stake, error) {
	var stakes []models.Stake

	err := repository.db.Model(&stakes).
		Column("Coin", "Validator", "OwnerAddress._").
		Where("owner_address.address = ?", address).
		Order("bip_value DESC").
		Select()

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

// Get paginated list of stakes by validator
func (repository Repository) GetPaginatedByValidator(
	validator models.Validator,
	pagination *tools.Pagination,
) ([]models.Stake, error) {
	var stakes []models.Stake
	var err error

	pagination.Total, err = repository.db.Model(&stakes).
		Column("Coin", "OwnerAddress.address").
		Where("validator_id = ?", validator.ID).
		Order("bip_value DESC").
		Apply(pagination.Filter).
		SelectAndCount()

	return stakes, err
}

func (repository Repository) GetMinStakes() ([]models.Stake, error) {
	var stakes []models.Stake

	err := repository.db.Model(&stakes).
		ColumnExpr("min(bip_value) as bip_value").
		Column("validator_id").
		Group("validator_id").
		Select()

	return stakes, err
}
