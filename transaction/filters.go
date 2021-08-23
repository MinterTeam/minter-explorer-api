package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/go-pg/pg/v10/orm"
	"strconv"
)

// TODO: replace string in StartBlock, EndBlock to int
type BlockFilter struct {
	BlockId uint64
}

func (f BlockFilter) Filter(q *orm.Query) (*orm.Query, error) {
	q = q.Where("transaction.block_id = ?", f.BlockId)

	return q, nil
}

// TODO: replace string in StartBlock, EndBlock to int
type ValidatorFilter struct {
	ValidatorPubKey string
	StartBlock      *string
	EndBlock        *string
}

func (f ValidatorFilter) Filter(q *orm.Query) (*orm.Query, error) {
	q = q.Join("LEFT JOIN transaction_validator").
		JoinOn("transaction_validator.transaction_id = transaction.id").
		Join("LEFT JOIN validators").
		JoinOn("validators.id = transaction_validator.validator_id").
		Where("validators.public_key = ?", f.ValidatorPubKey)

	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}

const (
	SendTypeIncoming  = "incoming"
	SendTypeOutcoming = "outcoming"
)

type SelectFilter struct {
	SendType   *string
	StartBlock *string
	EndBlock   *string
}

func (f SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.SendType != nil && *f.SendType == SendTypeIncoming {
		q.Where("transaction.from_address_id != a.id")
	}

	if f.SendType != nil && *f.SendType == SendTypeOutcoming {
		q.Where("transaction.from_address_id = a.id")
	}

	blocksRange := blocks.RangeSelectFilter{Prefix: "transaction.", StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}

type PoolsFilter struct {
	Coin0      string
	Coin1      string
	Token      string
	StartBlock *string
	EndBlock   *string
}

func (f PoolsFilter) Filter(q *orm.Query) (*orm.Query, error) {
	q = q.Join("LEFT JOIN transaction_liquidity_pool").
		JoinOn("transaction_liquidity_pool.transaction_id = transaction.id")

	if len(f.Token) > 0 {
		if id, err := strconv.Atoi(f.Token); err == nil {
			q = q.Where("transaction_liquidity_pool.liquidity_pool_id = ?", id)
		} else {
			q = q.Where("transaction_liquidity_pool.liquidity_pool_id = ?", helpers.GetPoolIdFromToken(f.Token))
		}
	}

	if len(f.Coin0) > 0 && len(f.Coin1) > 0 {
		q = q.Join("LEFT JOIN liquidity_pools").
			JoinOn("liquidity_pools.id = transaction_liquidity_pool.liquidity_pool_id")

		// filter by coins
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin0)
		symbol1, version1 := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin1)

		isSymbols := false
		_, err0 := strconv.Atoi(f.Coin0)
		_, err1 := strconv.Atoi(f.Coin1)
		if err0 != nil || err1 != nil {
			isSymbols = true
		}

		if isSymbols {
			q = q.Join("LEFT JOIN coins as first_coin").JoinOn("first_coin.id = liquidity_pools.first_coin_id").
				Where(`first_coin.symbol = ?`, symbol).Where(`first_coin.version = ?`, version)

			q = q.Join("LEFT JOIN coins as second_coin").JoinOn("second_coin.id = liquidity_pools.second_coin_id").
				Where(`second_coin.symbol = ?`, symbol1).Where(`second_coin.version = ?`, version1)
		}

		firstCoinAlias := "first_coin"
		secondCoinAlias := "second_coin"

		q = q.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
			if isSymbols {
				query = query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
						return query.Where(`"`+firstCoinAlias+`"."symbol" = ?`, symbol).
							Where(`"`+firstCoinAlias+`"."version" = ?`, version).
							Where(`"`+secondCoinAlias+`"."symbol" = ?`, symbol1).
							Where(`"`+secondCoinAlias+`"."version" = ?`, version1), nil
					}), nil
				})
			}

			if !isSymbols {
				query = query.WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.Where("liquidity_pools.first_coin_id = ?", f.Coin0).Where("liquidity_pools.second_coin_id = ?", f.Coin1), nil
				})
			}

			return query, nil
		}).WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
			if isSymbols {
				query = query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.WhereGroup(func(query *orm.Query) (*orm.Query, error) {
						return query.Where(`"`+firstCoinAlias+`"."symbol" = ?`, symbol1).
							Where(`"`+firstCoinAlias+`"."version" = ?`, version1).
							Where(`"`+secondCoinAlias+`"."symbol" = ?`, symbol).
							Where(`"`+secondCoinAlias+`"."version" = ?`, version), nil
					}), nil
				})
			}

			if !isSymbols {
				query = query.WhereOrGroup(func(query *orm.Query) (*orm.Query, error) {
					return query.Where("liquidity_pools.first_coin_id = ?", f.Coin1).Where("liquidity_pools.second_coin_id = ?", f.Coin0), nil
				})
			}

			return query, nil
		})
	}

	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}
