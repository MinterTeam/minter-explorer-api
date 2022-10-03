package validator

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
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

func (r Repository) GetByPublicKey(publicKey string) *models.Validator {
	var validator models.Validator

	err := r.db.Model(&validator).
		Relation("Stakes").
		Relation("OwnerAddress.address").
		Relation("RewardAddress.address").
		Relation("ControlAddress.address").
		Join("JOIN validator_public_keys ON validator_public_keys.validator_id = validator.id").
		Where("validator_public_keys.key = ?", publicKey).
		Select()

	if err != nil {
		return nil
	}

	return &validator
}

func (r Repository) GetTotalStakeByActiveValidators(ids []uint) string {
	var total string

	// get total stake of active validators
	err := r.db.Model((*models.Validator)(nil)).
		ColumnExpr("SUM(total_stake)").
		Where("id IN (?)", pg.In(ids)).
		Select(&total)

	helpers.CheckErr(err)

	return total
}

func (r Repository) GetActiveValidatorIds() []uint {
	var blockValidator models.BlockValidator
	var lastBlock models.Block
	var ids []uint

	lastBlockQuery := r.db.Model(&lastBlock).
		Column("id").
		Order("id DESC").
		Limit(1)

	// get active validators by last block
	err := r.db.Model(&blockValidator).
		Column("validator_id").
		Where("block_id = (?)", lastBlockQuery).
		Select(&ids)

	helpers.CheckErr(err)

	return ids
}

// Get active candidates count
func (r Repository) GetActiveCandidatesCount() int {
	var validator models.Validator

	count, err := r.db.Model(&validator).
		Where("status = ?", models.ValidatorStatusReady).
		Count()

	helpers.CheckErr(err)

	return count
}

// GetValidatorsAndStakes Get validators and stakes
func (r Repository) GetValidatorsAndStakes() []models.Validator {
	var validators []models.Validator

	err := r.db.Model(&validators).
		Relation("Stakes").
		Relation("OwnerAddress.address").
		Relation("RewardAddress.address").
		Relation("ControlAddress.address").
		Where("status is not null").
		Order("total_stake desc").
		Select()

	helpers.CheckErr(err)

	return validators
}

// GetValidators Get validators
func (r Repository) GetValidators() []models.Validator {
	var validators []models.Validator

	err := r.db.Model(&validators).
		Where("status is not null").
		Order("total_stake desc").
		Select()

	helpers.CheckErr(err)

	return validators
}

// Get validator bans
func (r Repository) GetBans(validator *models.Validator, pagination *tools.Pagination) (bans []models.ValidatorBan, err error) {
	pagination.Total, err = r.db.Model(&bans).
		Relation("Block").
		Where("validator_id = ?", validator.ID).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	return bans, err
}

// Get bans of validator list
func (r Repository) GetBansByValidatorIds(validatorIds []uint64, pagination *tools.Pagination) ([]models.ValidatorBan, error) {
	var bans []models.ValidatorBan

	err := r.db.Model(&bans).
		Relation("Block").
		Relation("Validator").
		Where(`validator_id in (?)`, pg.In(validatorIds)).
		Order("block_id DESC").
		Select()

	return bans, err
}
