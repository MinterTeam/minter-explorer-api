package validator

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/helpers"
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

func (repository Repository) GetByPublicKey(publicKey string) *models.Validator {
	var validator models.Validator

	err := repository.db.Model(&validator).
		Column("Stakes", "Stakes.Coin", "Stakes.OwnerAddress").
		Where("public_key = ?", publicKey).
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
	var ids []uint64

	// get active validators by last block
	err := repository.db.Model(&blockValidator).
		Column("validator_id").
		Where("block_id = ?", blocks.NewRepository(repository.db).GetLastBlock().ID).
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
