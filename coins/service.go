package coins

import (
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
)

type CoinService struct {
	DB *pg.DB
}

// Get list of coins
func (service *CoinService) GetList() *[]CoinResource {
	var coins []models.Coin
	err := service.DB.Model(&coins).Column("crr", "volume", "reserve_balance", "name", "symbol").Select()

	// check to existing row
	if err != nil {
		return nil
	}

	// transform models to resource
	var coinsResource []CoinResource
	for _, coin := range coins {
		coinsResource = append(coinsResource, TransformCoin(coin))
	}

	return &coinsResource
}

// Get coin detail by symbol
func (service *CoinService) GetBySymbol(symbol string) *CoinResource {
	coin := models.Coin{Symbol: symbol}

	// fetch data
	err := service.DB.Model(&coin).Where("symbol = ?", symbol).
		Column("crr", "volume", "reserve_balance", "name", "symbol").Select()

	// check to existing row
	if err != nil {
		return nil
	}

	// transform model to resource
	coinResource := TransformCoin(coin)

	return &coinResource
}