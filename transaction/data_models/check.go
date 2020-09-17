package data_models

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math/big"
)

type Check struct {
	RawCheck string
	Proof    string
	Check    CheckData
}

type CheckData struct {
	Coin     models.Coin
	GasCoin  models.Coin
	Nonce    []byte
	Value    *big.Int
	Sender   string
	DueBlock uint64
}
