package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/go-pg/pg/orm"
)

type SelectFilter struct {
	Addresses       []interface{}
	BlockId         *uint64
	StartBlock      *string
	EndBlock        *string
	ValidatorPubKey *string
}

func (f *SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if len(f.Addresses) > 0 {
		q = q.Join("LEFT OUTER JOIN transaction_outputs ON transaction_outputs.transaction_id = transaction.id").
			Join("JOIN addresses ON (addresses.id = transaction_outputs.to_address_id OR addresses.id = transaction.from_address_id)").
			WhereIn("addresses.address IN (?)", f.Addresses...)
	}

	if f.ValidatorPubKey != nil {
		q = q.Join("LEFT OUTER JOIN transaction_outputs ON transaction_outputs.transaction_id = transaction.id").
			Join("JOIN transaction_validator ON transaction_validator.transaction_id = transaction.id").
			Join("JOIN validators ON validators.public_key = ?", f.ValidatorPubKey)
	}

	if f.BlockId != nil {
		q = q.Where("transaction.block_id = ?", f.BlockId)
	}

	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}
