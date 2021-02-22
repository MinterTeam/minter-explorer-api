package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/go-pg/pg/v10/orm"
	"strconv"
)

type SelectPoolsFilter struct {
	Coin            *string
	ProviderAddress *string
}

func (f SelectPoolsFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.Coin != nil {
		if id, err := strconv.Atoi(*f.Coin); err == nil {
			q = q.Where("first_coin_id = ?", id).WhereOr("second_coin_id = ?", id)
		} else {
			symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(*f.Coin)
			q = q.WhereGroup(func(*orm.Query) (*orm.Query, error) {
				return q.Where(`"first_coin"."symbol" = ?`, symbol).Where(`"first_coin"."version" = ?`, version), nil
			}).WhereOrGroup(func(*orm.Query) (*orm.Query, error) {
				return q.Where(`"second_coin"."symbol" = ?`, symbol).Where(`"second_coin"."version" = ?`, version), nil
			})
		}
	}

	if f.ProviderAddress != nil {
		q = q.Join("LEFT OUTER JOIN address_liquidity_pools as provider").
			JoinOn("provider.liquidity_pool_id = liquidity_pool.id").
			Join("JOIN addresses as provider_address").
			JoinOn("provider_address.id = provider.address_id and provider_address.address = ?", helpers.RemoveMinterPrefix(*f.ProviderAddress))
	}

	return q, nil
}

type SelectByCoinsFilter struct {
	Coin0 string
	Coin1 string
	Token string
}

func (f SelectByCoinsFilter) Filter(tokenAlias, firstCoinAlias, secondCoinAlias string) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		// filter by token
		if len(f.Token) > 0 {
			if id, err := strconv.Atoi(f.Token); err == nil {
				q = q.Where("token_id = ?", id)
			} else {
				q = q.Where(`"`+tokenAlias+`"."symbol" = ?`, f.Token)
			}

			return q, nil
		}

		// filter by coins
		if id, err := strconv.Atoi(f.Coin0); err == nil {
			q = q.Where("first_coin_id = ?", id).WhereOr("second_coin_id = ?", id)
		} else {
			symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin0)
			q = q.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
				return query.Where(`"`+firstCoinAlias+`"."symbol" = ?`, symbol).Where(`"`+firstCoinAlias+`"."version" = ?`, version).
					WhereOr(`"`+secondCoinAlias+`"."symbol" = ?`, symbol).Where(`"`+secondCoinAlias+`"."version" = ?`, version), nil
			})
		}

		if id, err := strconv.Atoi(f.Coin1); err == nil {
			q = q.Where("second_coin_id = ?", id).WhereOr("first_coin_id = ?", id)
		} else {
			symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin1)
			q = q.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
				return query.Where(`"`+secondCoinAlias+`"."symbol" = ?`, symbol).Where(`"`+secondCoinAlias+`"."version" = ?`, version).
					WhereOr(`"`+firstCoinAlias+`"."symbol" = ?`, symbol).Where(`"`+secondCoinAlias+`"."version" = ?`, version), nil
			})
		}

		return q, nil
	}
}
