package events

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/go-pg/pg/v9/orm"
)

type SelectFilter struct {
	Address    string
	StartBlock *string
	EndBlock   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}

	return q.Where("address.address = ?", f.Address).Apply(blocksRange.Filter), nil
}
