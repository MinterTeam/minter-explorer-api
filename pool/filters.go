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
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin0)
		symbol1, version1 := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin1)

		isSymbols := false
		_, err0 := strconv.Atoi(f.Coin0)
		_, err1 := strconv.Atoi(f.Coin1)
		if err0 != nil || err1 != nil {
			isSymbols = true
		}

		q = q.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			query = query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
				return query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.Where(`"`+firstCoinAlias+`"."symbol" = ?`, symbol).
						Where(`"`+firstCoinAlias+`"."version" = ?`, version).
						Where(`"`+secondCoinAlias+`"."symbol" = ?`, symbol1).
						Where(`"`+secondCoinAlias+`"."version" = ?`, version1), nil
				}), nil
			})

			if !isSymbols {
				query = query.WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.Where("first_coin_id = ?", f.Coin0).Where("second_coin_id = ?", f.Coin1), nil
				})
			}

			return query, nil
		}).WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
			query = query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
				return query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.Where(`"`+firstCoinAlias+`"."symbol" = ?`, symbol1).
						Where(`"`+firstCoinAlias+`"."version" = ?`, version1).
						Where(`"`+secondCoinAlias+`"."symbol" = ?`, symbol).
						Where(`"`+secondCoinAlias+`"."version" = ?`, version), nil
				}), nil
			})

			if !isSymbols {
				query = query.WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.Where("first_coin_id = ?", f.Coin1).Where("second_coin_id = ?", f.Coin0), nil
				})
			}

			return query, nil
		})

		return q, nil
	}
}
