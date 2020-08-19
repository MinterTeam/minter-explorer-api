package data_resources

import (
	"encoding/base64"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type RedeemCheck struct {
	RawCheck string    `json:"raw_check"`
	Proof    string    `json:"proof"`
	Check    CheckData `json:"check"`
}

type CheckData struct {
	Coin     string `json:"coin"`
	GasCoin  string `json:"gas_coin"`
	Nonce    string `json:"nonce"`
	Value    string `json:"value"`
	Sender   string `json:"sender"`
	DueBlock uint64 `json:"due_block"`
}

func (RedeemCheck) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.RedeemCheckData)

	//TODO: handle error
	check, _ := TransformCheckData(data.GetRawCheck())

	return RedeemCheck{
		RawCheck: data.GetRawCheck(),
		Proof:    data.GetProof(),
		Check:    check,
	}
}

func TransformCheckData(raw string) (CheckData, error) {
	data, err := transaction.DecodeCheck(raw)
	if err != nil {
		return CheckData{}, err
	}

	sender, err := data.Sender()
	if err != nil {
		return CheckData{}, err
	}

	return CheckData{
		Coin:     string(data.Coin[:]),
		GasCoin:  string(data.GasCoin[:]),
		Nonce:    base64.StdEncoding.EncodeToString(data.Nonce),
		Value:    helpers.PipStr2Bip(data.Value.String()),
		Sender:   sender,
		DueBlock: data.DueBlock,
	}, nil
}
