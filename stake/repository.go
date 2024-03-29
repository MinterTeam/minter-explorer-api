package stake

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
	"math/big"
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
func (r Repository) GetAllByAddress(address string) ([]models.Stake, error) {
	var stakes []models.Stake

	err := r.db.Model(&stakes).
		Relation("Coin").
		Relation("Validator").
		Relation("OwnerAddress._").
		Where("owner_address.address = ?", address).
		Order("bip_value DESC").
		Select()

	return stakes, err
}

// Get total delegated bip value
func (r Repository) GetSumInBipValue() (string, error) {
	var sum string

	err := r.db.Model(&models.Stake{}).
		Where("is_kicked = false").
		ColumnExpr("SUM(bip_value)").
		Select(&sum)

	return sum, err
}

// Get total delegated sum by address
func (r Repository) GetSumInBipValueByAddress(address string) (string, error) {
	var sum string
	err := r.db.Model(&models.Stake{}).
		Relation("OwnerAddress._").
		ColumnExpr("SUM(bip_value)").
		Where("is_kicked = false").
		Where("owner_address.address = ?", address).
		Select(&sum)

	return sum, err
}

// Get paginated list of stakes by validator
func (r Repository) GetPaginatedByValidator(
	validator models.Validator,
	pagination *tools.Pagination,
) ([]models.Stake, error) {
	var stakes []models.Stake
	var err error

	pagination.Total, err = r.db.Model(&stakes).
		Relation("Coin").
		Relation("OwnerAddress.address").
		Where("validator_id = ?", validator.ID).
		Order("bip_value DESC").
		Apply(pagination.Filter).
		SelectAndCount()

	return stakes, err
}

func (r Repository) GetMinStakes() ([]models.Stake, error) {
	var stakes []models.Stake

	err := r.db.Model(&stakes).
		ColumnExpr("min(bip_value) as bip_value").
		Column("validator_id").
		Where("bip_value != 0").
		Where("is_kicked = false").
		Group("validator_id").
		Select()

	return stakes, err
}

func (r Repository) GetSumValueByCoin(coinID uint) (string, error) {
	var sum string

	err := r.db.Model(new(models.Stake)).
		ColumnExpr("SUM(value)").
		Where("coin_id = ?", coinID).
		Select(&sum)

	return sum, err
}

func (r Repository) GetDelegatorsCount() (count uint64, err error) {
	err = r.db.Model(new(models.Stake)).
		ColumnExpr("count (DISTINCT owner_address_id)").
		Select(&count)

	return count, err
}

func (r Repository) GetAddressValidatorIds(address string) (ids []uint64, err error) {
	err = r.db.Model(new(models.Stake)).
		Relation("OwnerAddress._").
		ColumnExpr("DISTINCT validator_id").
		Where("owner_address.address = ?", address).
		Select(&ids)

	return ids, err
}

func (r Repository) GetTotalStakeLocked() (*big.Int, error) {
	var tsl string

	err := r.db.Model(new(models.Stake)).
		ColumnExpr("sum(bip_value)").
		Where(`exists (select 1 from transactions where type = 37 and from_address_id = "stake"."owner_address_id")`).
		Select(&tsl)

	if err != nil {
		return nil, err
	}

	return helpers.StringToBigInt(tsl), err
}
