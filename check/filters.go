package check

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/go-pg/pg/v9/orm"
)

type SelectFilter struct {
	FromAddress string
	ToAddress   string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if len(f.FromAddress) != 0 {
		q = q.Where(`"from_address"."address" = ?`, helpers.RemoveMinterPrefix(f.FromAddress))
	}

	if len(f.ToAddress) != 0 {
		q = q.Where(`"to_address"."address" = ?`, helpers.RemoveMinterPrefix(f.ToAddress))
	}

	return q, nil
}
