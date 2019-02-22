package validator

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/models"
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

func (repository Repository) GetByPublicKey(publicKey string) models.Validator {
	var validator models.Validator

	err := repository.db.Model(&validator).
		Column("Stakes", "Stakes.Coin", "Stakes.OwnerAddress").
		Where("public_key = ?", publicKey).
		Select()

	helpers.CheckErr(err)

	return validator
}

func (repository Repository) GetTotalStake() string {
	var total string
	var blockValidator models.BlockValidator

	// get active validators by last block
	activeValidators := repository.db.Model(&blockValidator).
		Column("validator_id").
		Where("block_id = ?", blocks.NewRepository(repository.db).GetLastBlock().ID)

	// get total stake of active validators
	err := repository.db.Model((*models.Validator)(nil)).
		ColumnExpr("SUM(total_stake)").
		Where("id IN (?)", activeValidators).
		Select(&total)

	helpers.CheckErr(err)

	return total
}
