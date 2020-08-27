package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
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

func (repository Repository) GetByPublicKey(publicKey string) *models.Validator {
	var validator models.Validator

	err := repository.db.Model(&validator).
		Column("Stakes").
		Join("JOIN validator_public_keys ON validator_public_keys.validator_id = validator.id").
		Where("validator_public_keys.key = ?", publicKey).
		Select()

	if err != nil {
		return nil
	}

	return &validator
}

func (repository Repository) GetTotalStakeByActiveValidators(ids []uint64) string {
	var total string

	// get total stake of active validators
	err := repository.db.Model((*models.Validator)(nil)).
		ColumnExpr("SUM(total_stake)").
		Where("id IN (?)", pg.In(ids)).
		Select(&total)

	helpers.CheckErr(err)

	return total
}

func (repository Repository) GetActiveValidatorIds() []uint64 {
	var blockValidator models.BlockValidator
	var lastBlock models.Block
	var ids []uint64

	lastBlockQuery := repository.db.Model(&lastBlock).
		Column("id").
		Order("id DESC").
		Limit(1)

	// get active validators by last block
	err := repository.db.Model(&blockValidator).
		Column("validator_id").
		Where("block_id = (?)", lastBlockQuery).
		Select(&ids)

	helpers.CheckErr(err)

	return ids
}

// Get active candidates count
func (repository Repository) GetActiveCandidatesCount() int {
	var validator models.Validator

	count, err := repository.db.Model(&validator).
		Where("status = ?", models.ValidatorStatusReady).
		Count()

	helpers.CheckErr(err)

	return count
}

// Get validators
func (repository Repository) GetValidators() []models.Validator {
	var validators []models.Validator

	err := repository.db.Model(&validators).Select()
	helpers.CheckErr(err)

	return validators
}
