package events

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/go-pg/pg/orm"
)

type SelectFilter struct {
	Address    string
	StartBlock *string
	EndBlock   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	q = q.Where("address.address = ?", f.Address)

	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}

