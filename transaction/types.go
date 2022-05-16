package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Transaction interface {
	GetTags() map[string]string
	GetType() uint8
	GetData() json.RawMessage
	SetData(interface{})
	SetCommissionPriceCoin(c models.Coin)
	GetModel() interface{}
}

type ValidTx struct {
	tx *models.Transaction
}

func NewValidTx(tx *models.Transaction) *ValidTx {
	return &ValidTx{tx}
}

func (v *ValidTx) GetTags() map[string]string {
	return v.tx.Tags
}

func (v *ValidTx) GetType() uint8 {
	return v.tx.Type
}

func (v *ValidTx) GetData() json.RawMessage {
	return v.tx.Data
}

func (v *ValidTx) SetData(value interface{}) {
	v.tx.IData = value
}

func (v *ValidTx) SetCommissionPriceCoin(coin models.Coin) {
	v.tx.CommissionPriceCoin = coin
}

func (v *ValidTx) GetModel() interface{} {
	return v.tx
}

type invalidTx struct {
	tx *models.InvalidTransaction
}

func (i *invalidTx) GetTags() map[string]string {
	return i.tx.Tags
}

func (i *invalidTx) GetType() uint8 {
	return i.tx.Type
}

func (i *invalidTx) SetCommissionPriceCoin(coin models.Coin) {
	i.tx.CommissionPriceCoin = coin
}

func (i *invalidTx) GetData() json.RawMessage {
	return nil
}

func (i *invalidTx) SetData(value interface{}) {

}

func (i *invalidTx) GetModel() interface{} {
	return i.tx
}
