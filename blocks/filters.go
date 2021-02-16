package blocks

import (
	"github.com/go-pg/pg/v10/orm"
)

// TODO: replace string to int
type RangeSelectFilter struct {
	Prefix     string
	StartBlock *string
	EndBlock   *string
}

func (f RangeSelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartBlock != nil {
		q = q.Where(f.Prefix+"block_id >= ?", f.StartBlock)
	}

	if f.EndBlock != nil {
		q = q.Where(f.Prefix+"block_id <= ?", f.EndBlock)
	}

	return q, nil
}
