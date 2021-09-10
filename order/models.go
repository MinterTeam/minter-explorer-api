package order

import "github.com/MinterTeam/minter-explorer-extender/v2/models"

type OrderTransaction struct {
	tableName       struct{} `pg:"orders"`
	Id              uint64   `json:"id"`
	AddressId       uint64   `json:"address_id"`
	LiquidityPoolId uint64   `json:"liquidity_pool_id"`
	CoinSellId      uint64   `json:"coin_sell_id"      pg:",use_zero"`
	CoinSellVolume  string   `json:"coin_sell_volume"`
	CoinBuyId       uint64   `json:"coin_buy_id"       pg:",use_zero"`
	CoinBuyVolume   string   `json:"coin_buy_volume"`
	CreatedAtBlock  uint64   `json:"created_at_block"`
	IsCanceled      bool
	Status          models.OrderType      `json:"status"`
	Address         *models.Address       `json:"address"           pg:"rel:has-one,fk:address_id"`
	LiquidityPool   *models.LiquidityPool `json:"liquidity_pool"    pg:"rel:has-one,fk:liquidity_pool_id"`
	CoinSell        *models.Coin          `json:"coin_sell"         pg:"rel:has-one,fk:coin_sell_id"`
	CoinBuy         *models.Coin          `json:"coin_buy"          pg:"rel:has-one,fk:coin_buy_id"`
	Transaction     *models.Transaction   `json:"transaction"`
}
