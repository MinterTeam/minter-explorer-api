package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction/data"
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

func (Resource) Transform(model resource.ItemInterface) resource.Interface {
	tx := model.(models.Transaction)

	return Resource{
		Txn:       tx.ID,
		Hash:      tx.GetHash(),
		Nonce:     tx.Nonce,
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt,
		Fee:       helpers.Fee2Bip(tx.GetFee()),
		Type:      GetTypeAsText(tx.Type),
		Payload:   string(tx.Payload[:]),
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
	models.TxTypeSend:                {Model: new(models.SendTxData), Resource: data.SendResource{}, TypeText: "send"},
	models.TxTypeSellCoin:            {Model: new(models.SellCoinTxData), Resource: data.SellCoinResource{}, TypeText: "sellCoin"},
	models.TxTypeSellAllCoin:         {Model: new(models.SellAllCoinTxData), Resource: data.SellAllCoinResource{}, TypeText: "sellAllCoin"},
	models.TxTypeBuyCoin:             {Model: new(models.BuyCoinTxData), Resource: data.BuyCoinResource{}, TypeText: "buyCoin"},
	models.TxTypeCreateCoin:          {Model: new(models.CreateCoinTxData), Resource: data.CreateCoinResource{}, TypeText: "createCoin"},
	models.TxTypeDeclareCandidacy:    {Model: new(models.DeclareCandidacyTxData), Resource: data.DeclareCandidacyResource{}, TypeText: "declareCandidacy"},
	models.TxTypeDelegate:            {Model: new(models.DelegateTxData), Resource: data.DelegateResource{}, TypeText: "delegate"},
	models.TxTypeUnbound:             {Model: new(models.UnbondTxData), Resource: data.UnbondResource{}, TypeText: "unbond"},
	models.TxTypeRedeemCheck:         {Model: new(models.RedeemCheckTxData), Resource: data.RedeemCheckResource{}, TypeText: "redeemCheckData"},
	models.TxTypeMultiSig:            {Model: new(models.CreateMultisigTxData), Resource: data.CreateMultisigResource{}, TypeText: "multiSig"},
	models.TxTypeMultiSend:           {Model: new(models.TransactionOutput), Resource: data.MultisendResource{}, TypeText: "multiSend"},
	models.TxTypeEditCandidate:       {Model: new(models.EditCandidateTxData), Resource: data.EditCandidateResource{}, TypeText: "editCandidate"},
	models.TxTypeSetCandidateOnline:  {Model: new(models.SetCandidateTxData), Resource: data.SetCandidateResource{}, TypeText: "setCandidateOnData"},
	models.TxTypeSetCandidateOffline: {Model: new(models.SetCandidateTxData), Resource: data.SetCandidateResource{}, TypeText: "setCandidateOffData"},
}

func TransformTxData(tx models.Transaction) resource.Interface {
	config := transformConfig[tx.Type]

	var val resource.ItemInterface
	if tx.Type == models.TxTypeMultiSend {
		val = tx.TxOutput
	} else {
		val = reflect.New(reflect.TypeOf(config.Model).Elem()).Interface()
		err := json.Unmarshal(tx.Data, val)
		helpers.CheckErr(err)
	}

	return config.Resource.Transform(val)
}

func GetTypeAsText(txType uint8) string {
	return transformConfig[txType].TypeText
}
