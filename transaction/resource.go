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
		Data:      TransformTxData(tx.Type, tx.Data),
	}
}

func TransformTxData(txType uint8, txData json.RawMessage) resource.Interface {
	switch txType {
	case models.TxTypeSend:
		return TransformTxDataModelToResource(txData, new(models.SendTxData), data.SendResource{})
	case models.TxTypeSellCoin:
		return TransformTxDataModelToResource(txData, new(models.SellCoinTxData), data.SellCoinResource{})
	case models.TxTypeSellAllCoin:
		return TransformTxDataModelToResource(txData, new(models.SellAllCoinTxData), data.SellAllCoinResource{})
	case models.TxTypeBuyCoin:
		return TransformTxDataModelToResource(txData, new(models.BuyCoinTxData), data.BuyCoinResource{})
	case models.TxTypeCreateCoin:
		return TransformTxDataModelToResource(txData, new(models.CreateCoinTxData), data.CreateCoinResource{})
	case models.TxTypeDeclareCandidacy:
		return TransformTxDataModelToResource(txData, new(models.DeclareCandidacyTxData), data.DeclareCandidacyResource{})
	case models.TxTypeDelegate:
		return TransformTxDataModelToResource(txData, new(models.DelegateTxData), data.DelegateResource{})
	case models.TxTypeUnbound:
		return TransformTxDataModelToResource(txData, new(models.UnbondTxData), data.UnbondResource{})
	case models.TxTypeRedeemCheck:
		return TransformTxDataModelToResource(txData, new(models.RedeemCheckTxData), data.RedeemCheckResource{})
	case models.TxTypeMultiSig:
		return TransformTxDataModelToResource(txData, new(models.CreateMultisigTxData), data.CreateMultisigResource{})
	case models.TxTypeMultiSend:
		return TransformTxDataModelToResource(txData, new(models.MultiSendTxData), data.MultisendResource{})
	case models.TxTypeEditCandidate:
		return TransformTxDataModelToResource(txData, new(models.EditCandidateTxData), data.EditCandidateResource{})
	case models.TxTypeSetCandidateOnline, models.TxTypeSetCandidateOffline:
		return TransformTxDataModelToResource(txData, new(models.SetCandidateTxData), data.SetCandidateResource{})
	}

	return nil
}

func GetTypeAsText(txType uint8) string {
	var type2text = map[uint8]string {
		models.TxTypeSend: "send",
		models.TxTypeSellCoin: "sellCoin",
		models.TxTypeSellAllCoin: "sellAllCoin",
		models.TxTypeBuyCoin: "buyCoin",
		models.TxTypeCreateCoin: "createCoin",
		models.TxTypeDeclareCandidacy: "declareCandidacy",
		models.TxTypeDelegate: "delegate",
		models.TxTypeUnbound: "unbond",
		models.TxTypeRedeemCheck: "redeemCheckData",
		models.TxTypeMultiSig: "multiSig",
		models.TxTypeMultiSend: "multiSend",
		models.TxTypeEditCandidate: "editCandidate",
		models.TxTypeSetCandidateOnline: "setCandidateOnData",
		models.TxTypeSetCandidateOffline: "setCandidateOffData",
	}

	return type2text[txType]
}

func TransformTxDataModelToResource(raw json.RawMessage, model resource.ItemInterface, resource resource.Interface) resource.Interface {
	val := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	err := json.Unmarshal(raw, val)
	helpers.CheckErr(err)

	return resource.Transform(val)
}
