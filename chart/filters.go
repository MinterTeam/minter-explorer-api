package chart

import "github.com/go-pg/pg/orm"

type SelectFilter struct {
	Scale     string
	StartTime *string
	EndTime   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.StartTime != nil {
		q = q.Where("created_at >= ?", f.StartTime)
	}

	if f.EndTime != nil {
		q = q.Where("created_at <= ?", f.EndTime)
	}

	return q.ColumnExpr("date_trunc(?, created_at) as time", f.Scale).Group("time").Order("time"), nil
}
