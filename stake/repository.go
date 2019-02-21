package stake

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

func (repository Repository) GetByAddress(address string) []models.Stake {
	var stakes []models.Stake

	err := repository.db.Model(&stakes).Column("Coin.symbol", "Validator.public_key").Select()
	helpers.CheckErr(err)

	return stakes
}
