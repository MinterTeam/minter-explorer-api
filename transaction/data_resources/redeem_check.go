package data_resources

import (
	"encoding/base64"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/MinterTeam/minter-go-node/core/check"
)

type RedeemCheck struct {
	RawCheck string    `json:"raw_check"`
	Proof    string    `json:"proof"`
	Check    CheckData `json:"check"`
}

type CheckData struct {
	Coin     string `json:"coin"`
	Nonce    string `json:"nonce"`
	Value    string `json:"value"`
	Sender   string `json:"sender"`
	DueBlock uint64 `json:"due_block"`
}

func (RedeemCheck) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.RedeemCheckTxData)

	return RedeemCheck{
		RawCheck: data.RawCheck,
		Proof:    data.Proof,
		Check:    TransformCheckData(data.RawCheck),
	}
}

func TransformCheckData(raw string) CheckData {
	decoded, err := base64.StdEncoding.DecodeString(raw)
	helpers.CheckErr(err)

	data, err := check.DecodeFromBytes(decoded)
	helpers.CheckErr(err)

	sender, err := data.Sender()
	helpers.CheckErr(err)

	return CheckData{
		Coin:     data.Coin.String(),
		Nonce:    base64.StdEncoding.EncodeToString(data.Nonce[:]),
		Value:    helpers.PipStr2Bip(data.Value.String()),
		Sender:   sender.String(),
		DueBlock: data.DueBlock,
	}
}
