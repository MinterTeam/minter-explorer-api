package tools

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/go-pg/pg/orm"
)

type BlocksRangeFilter struct {
	StartBlock *string
	EndBlock   *string
}

func (f BlocksRangeFilter) Apply(q *orm.Query) (*orm.Query, error) {
	if f.StartBlock != nil {
		q = q.Where("block_id >= ?", f.StartBlock)
	}

	if f.EndBlock != nil {
		q = q.Where("block_id <= ?", f.EndBlock)
	}

	return q, nil
}

type EventsFilter struct {
	Address    string
	StartBlock *string
	EndBlock   *string
}

func (f EventsFilter) Apply(q *orm.Query) (*orm.Query, error) {
	var err error

	q = q.Where("address.address = ?", f.Address)

	blockRangeFilter := BlocksRangeFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q, err = blockRangeFilter.Apply(q)
	helpers.CheckErr(err)

	return q, nil
}
