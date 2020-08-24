package waitlist

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
)

type Repository struct {
	db *pg.DB
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetListByAddress(address string) ([]models.StakeKick, error) {
	var wl []models.StakeKick

	err := r.db.Model(&wl).
		Column("Coin", "Validator").
		Join("addresses").
		Where("address_id = ?", address).
		Select()

	if err != nil {
		return nil, err
	}

	return wl, nil
}
