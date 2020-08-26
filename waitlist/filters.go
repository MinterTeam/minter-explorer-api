package waitlist

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/go-pg/pg/v9/orm"
)

type SelectFilter struct {
	PublicKey *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.PublicKey != nil {
		q = q.Where("validator.public_key = ?", helpers.RemoveMinterPrefix(*f.PublicKey))
	}

	return q, nil
}
