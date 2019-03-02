package transaction

import (
	"encoding/base64"
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction/data_resources"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"reflect"
	"time"
)

type Resource struct {
	Txn       uint64                 `json:"txn"`
	Hash      string                 `json:"hash"`
	Nonce     uint64                 `json:"nonce"`
	Block     uint64                 `json:"block"`
	Timestamp time.Time              `json:"timestamp"`
	Fee       string                 `json:"fee"`
	Type      string                 `json:"type"`
	Payload   string                 `json:"payload"`
	From      string                 `json:"from"`
	Data      resource.ItemInterface `json:"data"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	tx := model.(models.Transaction)

	return Resource{
		Txn:       tx.ID,
		Hash:      tx.GetHash(),
		Nonce:     tx.Nonce,
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt,
		Fee:       helpers.Fee2Bip(tx.GetFee()),
		Type:      GetTypeAsText(tx.Type),
		Payload:   base64.StdEncoding.EncodeToString(tx.Payload[:]),
		From:      tx.FromAddress.GetAddress(),
		Data:      TransformTxData(tx),
	}
}

type TransformTxConfig struct {
	Model    resource.ItemInterface
	Resource resource.Interface
	TypeText string
}

var transformConfig = map[uint8]TransformTxConfig{
	models.TxTypeSend:                {Model: new(models.SendTxData), Resource: data_resources.Send{}, TypeText: "send"},
	models.TxTypeSellCoin:            {Model: new(models.SellCoinTxData), Resource: data_resources.SellCoin{}, TypeText: "sellCoin"},
	models.TxTypeSellAllCoin:         {Model: new(models.SellAllCoinTxData), Resource: data_resources.SellAllCoin{}, TypeText: "sellAllCoin"},
	models.TxTypeBuyCoin:             {Model: new(models.BuyCoinTxData), Resource: data_resources.BuyCoin{}, TypeText: "buyCoin"},
	models.TxTypeCreateCoin:          {Model: new(models.CreateCoinTxData), Resource: data_resources.CreateCoin{}, TypeText: "createCoin"},
	models.TxTypeDeclareCandidacy:    {Model: new(models.DeclareCandidacyTxData), Resource: data_resources.DeclareCandidacy{}, TypeText: "declareCandidacy"},
	models.TxTypeDelegate:            {Model: new(models.DelegateTxData), Resource: data_resources.Delegate{}, TypeText: "delegate"},
	models.TxTypeUnbound:             {Model: new(models.UnbondTxData), Resource: data_resources.Unbond{}, TypeText: "unbond"},
	models.TxTypeRedeemCheck:         {Model: new(models.RedeemCheckTxData), Resource: data_resources.RedeemCheck{}, TypeText: "redeemCheckData"},
	models.TxTypeMultiSig:            {Model: new(models.CreateMultisigTxData), Resource: data_resources.CreateMultisig{}, TypeText: "multiSig"},
	models.TxTypeMultiSend:           {Model: new(models.MultiSendTxData), Resource: data_resources.Multisend{}, TypeText: "multiSend"},
	models.TxTypeEditCandidate:       {Model: new(models.EditCandidateTxData), Resource: data_resources.EditCandidate{}, TypeText: "editCandidate"},
	models.TxTypeSetCandidateOnline:  {Model: new(models.SetCandidateTxData), Resource: data_resources.SetCandidate{}, TypeText: "setCandidateOnData"},
	models.TxTypeSetCandidateOffline: {Model: new(models.SetCandidateTxData), Resource: data_resources.SetCandidate{}, TypeText: "setCandidateOffData"},
}

func TransformTxData(tx models.Transaction) resource.Interface {
	config := transformConfig[tx.Type]

	val := reflect.New(reflect.TypeOf(config.Model).Elem()).Interface()
	err := json.Unmarshal(tx.Data, val)
	helpers.CheckErr(err)

	return config.Resource.Transform(val, tx)
}

func GetTypeAsText(txType uint8) string {
	return transformConfig[txType].TypeText
}
