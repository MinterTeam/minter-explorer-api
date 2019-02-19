package blocks

import (
	"github.com/go-pg/pg/orm"
)

// TODO: replace string to int
type RangeSelectFilter struct {
	StartBlock *string
	EndBlock   *string
}

func (f RangeSelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartBlock != nil {
		q = q.Where("block_id >= ?", f.StartBlock)
	}

	if f.EndBlock != nil {
		q = q.Where("block_id <= ?", f.EndBlock)
	}

	return q, nil
}
