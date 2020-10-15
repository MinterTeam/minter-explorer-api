package coins

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
)

type Repository struct {
	DB        *pg.DB
	baseModel *models.Coin
}

var GlobalRepository *Repository

func NewRepository(db *pg.DB) *Repository {
	GlobalRepository = &Repository{
		DB: db,
	}

	return GlobalRepository
}

// Get list of coins
func (repository *Repository) GetCoins() []models.Coin {
	var coins []models.Coin

	err := repository.DB.Model(&coins).
		Column("OwnerAddress").
		Where("deleted_at IS NULL").
		Order("reserve DESC").
		Select()

	helpers.CheckErr(err)

	return coins
}

// Get coin detail by symbol
func (repository *Repository) GetBySymbolAndVersion(symbol string, version *uint64) []models.Coin {
	var coins []models.Coin

	query := repository.DB.Model(&coins).
		Column("OwnerAddress").
		Where("symbol LIKE ?", fmt.Sprintf("%%%s%%", symbol)).
		Where("deleted_at IS NULL").
		Order("reserve DESC")

	if version != nil {
		query.Where("version = ?", version)
	}

	err := query.Select()
	helpers.CheckErr(err)

	return coins
}

type CustomCoinsStatusData struct {
	ReserveSum string
	Count      uint
}

// Get custom coins data for status page
func (repository *Repository) GetCustomCoinsStatusData() (CustomCoinsStatusData, error) {
	var data CustomCoinsStatusData

	err := repository.DB.
		Model(&models.Coin{}).
		ColumnExpr("SUM(reserve) as reserve_sum, COUNT(*) as count").
		Where("id != ?", 0).
		Select(&data)

	return data, err
}

func (repository *Repository) FindByID(id uint) (models.Coin, error) {
	var coin models.Coin

	if id == 0 && repository.baseModel != nil {
		return *repository.baseModel, nil
	}

	err := repository.DB.Model(&coin).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Select()

	if id == 0 && repository.baseModel == nil {
		repository.baseModel = &coin
	}

	return coin, err
}
