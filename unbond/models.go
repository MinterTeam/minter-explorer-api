package unbond

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"time"
)

type UnbondMoveStake struct {
	BlockId       uint
	CoinId        uint
	ValidatorId   uint
	ToValidatorId *uint
	Value         string
	CreatedAt     time.Time
	MinterAddress string
	Address       *models.Address
	Coin          *models.Coin
	FromValidator *models.Validator
	ToValidator   *models.Validator
}
