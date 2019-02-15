package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"strconv"
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
		Fee:       strconv.FormatUint(tx.GetFee(), 10),
		Type:      getTypeAsText(tx.Type),
		Status:    "success",
		Payload:   string(tx.Payload[:]),
		From:      tx.FromAddress.GetAddress(),
		Data:      TransformTxData(tx.Type, tx.Data),
	}
}

func TransformTxData(txType uint8, txData json.RawMessage) resource.Interface {
	switch txType {
	case models.TxTypeSend:
		return new(SendDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeSellCoin:
		return new(SellCoinDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeSellAllCoin:
		return new(SellAllCoinDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeBuyCoin:
		return new(BuyCoinDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeCreateCoin:
		return new(CreateCoinDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeDeclareCandidacy:
		return new(DeclareCandidacyDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeDelegate:
		return new(DelegateDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeUnbound:
		return new(UnbondDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeRedeemCheck:
		return new(RedeemCheckDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeMultiSig:
		return new(CreateMultisigDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeMultiSend:
		return new(MultisendDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeEditCandidate:
		return new(EditCandidateDataResource).TransformFromJsonRaw(txData)
	case models.TxTypeSetCandidateOnline:
		return new(SetCandidateOnResource).TransformFromJsonRaw(txData)
	case models.TxTypeSetCandidateOffline:
		return new(SetCandidateOffDataResource).TransformFromJsonRaw(txData)
	}

	return nil
}

func getTypeAsText(txType uint8) string {
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
