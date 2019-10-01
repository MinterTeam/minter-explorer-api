package aggregated_reward

import "github.com/go-pg/pg/orm"

type SelectFilter struct {
	Address   string
	StartTime *string
	EndTime   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartTime != nil {
		q = q.Where("time_id >= ?", f.StartTime)
	}

	if f.EndTime != nil {
		q = q.Where("time_id <= ?", f.EndTime)
	}

	return q.Where("address.address = ?", f.Address), nil
}
