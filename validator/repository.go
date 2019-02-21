package validator

import (
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

	err := repository.db.Model((*models.Validator)(nil)).ColumnExpr("SUM(total_stake)").Select(&total)
	helpers.CheckErr(err)

	return total
}
