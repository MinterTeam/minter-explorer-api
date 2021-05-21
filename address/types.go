package address

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math/big"
)

type Balance struct {
	Model              *models.Address
	TotalBalanceSum    *big.Int
	TotalBalanceSumUSD *big.Float
	StakeBalanceSum    *big.Int
	StakeBalanceSumUSD *big.Float
}
