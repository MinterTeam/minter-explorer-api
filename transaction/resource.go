package transaction

import (
	"encoding/base64"
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction/data_resources"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"reflect"
	"strconv"
	"time"
)

type Resource struct {
	Txn       uint64                 `json:"txn"`
	Hash      string                 `json:"hash"`
	Nonce     uint64                 `json:"nonce"`
	Block     uint64                 `json:"height"`
	Timestamp string                 `json:"timestamp"`
	GasCoin   string                 `json:"gas_coin"`
	Gas       string                 `json:"gas"`
	GasPrice  uint64                 `json:"gas_price"`
	Fee       string                 `json:"fee"`
	Type      uint8                  `json:"type"`
	Payload   string                 `json:"payload"`
	From      string                 `json:"from"`
	Data      resource.ItemInterface `json:"data"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := model.(models.Transaction)

	return Resource{
		Txn:       tx.ID,
		Hash:      tx.GetHash(),
		Nonce:     tx.Nonce,
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt.Format(time.RFC3339),
		Gas:       strconv.FormatUint(tx.Gas, 10),
		GasPrice:  tx.GasPrice,
		Fee:       helpers.Fee2Bip(tx.GetFee()),
		GasCoin:   tx.GasCoin.Symbol,
		Type:      tx.Type,
		Payload:   base64.StdEncoding.EncodeToString(tx.Payload[:]),
		From:      tx.FromAddress.GetAddress(),
		Data:      TransformTxData(tx),
	}
}

type TransformTxConfig struct {
	Model    resource.ItemInterface
	Resource resource.Interface
}

var transformConfig = map[uint8]TransformTxConfig{
	uint8(transaction.TypeSend):                {Model: new(models.SendTxData), Resource: data_resources.Send{}},
	uint8(transaction.TypeSellCoin):            {Model: new(models.SellCoinTxData), Resource: data_resources.SellCoin{}},
	uint8(transaction.TypeSellAllCoin):         {Model: new(models.SellAllCoinTxData), Resource: data_resources.SellAllCoin{}},
	uint8(transaction.TypeBuyCoin):             {Model: new(models.BuyCoinTxData), Resource: data_resources.BuyCoin{}},
	uint8(transaction.TypeCreateCoin):          {Model: new(models.CreateCoinTxData), Resource: data_resources.CreateCoin{}},
	uint8(transaction.TypeDeclareCandidacy):    {Model: new(models.DeclareCandidacyTxData), Resource: data_resources.DeclareCandidacy{}},
	uint8(transaction.TypeDelegate):            {Model: new(models.DelegateTxData), Resource: data_resources.Delegate{}},
	uint8(transaction.TypeUnbond):              {Model: new(models.UnbondTxData), Resource: data_resources.Unbond{}},
	uint8(transaction.TypeRedeemCheck):         {Model: new(models.RedeemCheckTxData), Resource: data_resources.RedeemCheck{}},
	uint8(transaction.TypeCreateMultisig):      {Model: new(models.CreateMultisigTxData), Resource: data_resources.CreateMultisig{}},
	uint8(transaction.TypeMultisend):           {Model: new(models.MultiSendTxData), Resource: data_resources.Multisend{}},
	uint8(transaction.TypeEditCandidate):       {Model: new(models.EditCandidateTxData), Resource: data_resources.EditCandidate{}},
	uint8(transaction.TypeSetCandidateOnline):  {Model: new(models.SetCandidateTxData), Resource: data_resources.SetCandidate{}},
	uint8(transaction.TypeSetCandidateOffline): {Model: new(models.SetCandidateTxData), Resource: data_resources.SetCandidate{}},
}

func TransformTxData(tx models.Transaction) resource.Interface {
	config := transformConfig[tx.Type]

	val := reflect.New(reflect.TypeOf(config.Model).Elem()).Interface()
	err := json.Unmarshal(tx.Data, val)
	helpers.CheckErr(err)

	return config.Resource.Transform(val, tx)
}
