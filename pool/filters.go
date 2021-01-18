package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/go-pg/pg/v9/orm"
)

type SelectPoolsFilter struct {
	CoinId          *uint64
	ProviderAddress *string
}

func (f SelectPoolsFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.CoinId != nil {
		q = q.Where("first_coin_id = ?", f.CoinId).WhereOr("second_coin_id = ?", f.CoinId)
	}

	if f.ProviderAddress != nil {
		q = q.Join("LEFT OUTER JOIN address_liquidity_pools as provider").
			JoinOn("provider.liquidity_pool_id = liquidity_pool.id").
			Join("JOIN addresses as provider_address").
			JoinOn("provider_address.id = provider.address_id and provider_address.address = ?", helpers.RemoveMinterPrefix(*f.ProviderAddress))
	}

	return q, nil
}
