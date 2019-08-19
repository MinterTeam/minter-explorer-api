package reward

import (
	"github.com/go-pg/pg/orm"
)

type AggregatedSelectFilter struct {
	Address    string
	StartBlock *string
	EndBlock   *string
}

func (f AggregatedSelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartBlock != nil {
		q = q.Where("from_block_id >= ?", f.StartBlock)
	}

	if f.EndBlock != nil {
		q = q.Where("to_block_id <= ?", f.EndBlock)
	}

	return q.Where("address.address = ?", f.Address), nil
}

