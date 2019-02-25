package chart

import "github.com/go-pg/pg/orm"

type SelectFilter struct {
	Scale     string
	StartTime *string
	EndTime   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartTime != nil {
		q = q.Where("block.created_at >= ?", f.StartTime)
	}

	if f.EndTime != nil {
		q = q.Where("block.created_at <= ?", f.EndTime)
	}

	return q.Column("Block._").ColumnExpr("date_trunc(?, block.created_at) as time", f.Scale).Group("time").Order("time"), nil
}
