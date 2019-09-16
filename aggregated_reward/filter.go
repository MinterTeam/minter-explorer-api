package aggregated_reward

import "github.com/go-pg/pg/orm"

type SelectFilter struct {
	Address    string
	StartBlock *string
	EndBlock   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartBlock != nil {
		q = q.Where("from_block_id >= ?", f.StartBlock)
	}

	if f.EndBlock != nil {
		q = q.Where("to_block_id <= ?", f.EndBlock)
	}

	return q.Where("address.address = ?", f.Address), nil
}
