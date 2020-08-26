package waitlist

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetListByAddress(address string, filter SelectFilter) ([]models.StakeKick, error) {
	var wl []models.StakeKick

	err := r.db.Model(&wl).
		Column("Coin", "Validator").
		Join("JOIN addresses ON addresses.id = stake_kick.address_id").
		Where("addresses.address = ?", address).
		Apply(filter.Filter).
		Select()

	if err != nil {
		return nil, err
	}

	return wl, nil
}
