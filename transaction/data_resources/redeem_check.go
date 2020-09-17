package data_resources

import (
	"encoding/base64"
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	txModels "github.com/MinterTeam/minter-explorer-api/transaction/data_models"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type RedeemCheck struct {
	RawCheck string    `json:"raw_check"`
	Proof    string    `json:"proof"`
	Check    CheckData `json:"check"`
}

type CheckData struct {
	Coin     resource.Interface `json:"coin"`
	GasCoin  resource.Interface `json:"gas_coin"`
	Nonce    string             `json:"nonce"`
	Value    string             `json:"value"`
	Sender   string             `json:"sender"`
	DueBlock uint64             `json:"due_block"`
}

func (RedeemCheck) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := params[0].(models.Transaction)
	model := tx.IData.(txModels.Check)

	return RedeemCheck{
		RawCheck: model.RawCheck,
		Proof:    model.Proof,
		Check: CheckData{
			Coin:     new(coins.IdResource).Transform(model.Check.Coin),
			GasCoin:  new(coins.IdResource).Transform(model.Check.GasCoin),
			Nonce:    base64.StdEncoding.EncodeToString(model.Check.Nonce),
			Value:    helpers.PipStr2Bip(model.Check.Value.String()),
			Sender:   model.Check.Sender,
			DueBlock: model.Check.DueBlock,
		},
	}
}
