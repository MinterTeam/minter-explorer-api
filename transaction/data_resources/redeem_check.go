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
	Data     CheckData `json:"data"`
}

type CheckData struct {
	Coin     string `json:"coin"`
	Nonce    uint64 `json:"nonce"`
	Value    string `json:"value"`
	Sender   string `json:"sender"`
	DueBlock uint64 `json:"due_block"`
}

func (RedeemCheck) Transform(txData resource.ItemInterface, params ...interface{}) resource.Interface {
	data := txData.(*models.RedeemCheckTxData)

	return RedeemCheck{
		RawCheck: data.RawCheck,
		Proof:    data.Proof,
		Data:     TransformCheckData(data.RawCheck),
	}
}

func TransformCheckData(raw string) CheckData {
	decoded, _ := base64.StdEncoding.DecodeString(raw)
	data, _ := check.DecodeFromBytes(decoded)
	sender, _ := data.Sender()

	return CheckData{
		Coin:     data.Coin.String(),
		Nonce:    data.Nonce,
		Value:    helpers.PipStr2Bip(data.Value.String()),
		Sender:   sender.String(),
		DueBlock: data.DueBlock,
	}
}
