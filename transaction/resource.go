package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction/data"
	"github.com/MinterTeam/minter-explorer-extender/models"
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
	Status    string                 `json:"status"`
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
		return new(data.SendResource).TransformFromJsonRaw(txData)
	case models.TxTypeSellCoin:
		return new(data.SellCoinResource).TransformFromJsonRaw(txData)
	case models.TxTypeSellAllCoin:
		return new(data.SellAllCoinResource).TransformFromJsonRaw(txData)
	case models.TxTypeBuyCoin:
		return new(data.BuyCoinResource).TransformFromJsonRaw(txData)
	case models.TxTypeCreateCoin:
		return new(data.CreateCoinResource).TransformFromJsonRaw(txData)
	case models.TxTypeDeclareCandidacy:
		return new(data.DeclareCandidacyResource).TransformFromJsonRaw(txData)
	case models.TxTypeDelegate:
		return new(data.DelegateResource).TransformFromJsonRaw(txData)
	case models.TxTypeUnbound:
		return new(data.UnbondResource).TransformFromJsonRaw(txData)
	case models.TxTypeRedeemCheck:
		return new(data.RedeemCheckResource).TransformFromJsonRaw(txData)
	case models.TxTypeMultiSig:
		return new(data.CreateMultisigResource).TransformFromJsonRaw(txData)
	case models.TxTypeMultiSend:
		return new(data.MultisendResource).TransformFromJsonRaw(txData)
	case models.TxTypeEditCandidate:
		return new(data.EditCandidateResource).TransformFromJsonRaw(txData)
	case models.TxTypeSetCandidateOnline:
		return new(data.SetCandidateOnResource).TransformFromJsonRaw(txData)
	case models.TxTypeSetCandidateOffline:
		return new(data.SetCandidateOffResource).TransformFromJsonRaw(txData)
	}

	return nil
}

func GetTypeAsText(txType uint8) string {
	switch txType {
	case models.TxTypeSend:
		return "send"
	case models.TxTypeSellCoin:
		return "sellCoin"
	case models.TxTypeSellAllCoin:
		return "sellAllCoin"
	case models.TxTypeBuyCoin:
		return "buyCoin"
	case models.TxTypeCreateCoin:
		return "createCoin"
	case models.TxTypeDeclareCandidacy:
		return "declareCandidacy"
	case models.TxTypeDelegate:
		return "delegate"
	case models.TxTypeUnbound:
		return "unbond"
	case models.TxTypeRedeemCheck:
		return "redeemCheckData"
	case models.TxTypeMultiSig:
		return "multiSig"
	case models.TxTypeMultiSend:
		return "multiSend"
	case models.TxTypeEditCandidate:
		return "editCandidate"
	case models.TxTypeSetCandidateOnline:
		return "setCandidateOnData"
	case models.TxTypeSetCandidateOffline:
		return "setCandidateOffData"
	}

	return ""
}
