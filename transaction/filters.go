package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/go-pg/pg/v9/orm"
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

		if id, err := strconv.Atoi(f.Coin0); err == nil {
			q = q.Where("liquidity_pools.first_coin_id = ?", id)
		} else {
			symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin0)
			q = q.Join("LEFT JOIN coins as first_coin").JoinOn("first_coin.id = liquidity_pools.first_coin_id").
				Where(`first_coin.symbol = ?`, symbol).Where(`first_coin.version = ?`, version)
		}

		if id, err := strconv.Atoi(f.Coin1); err == nil {
			q = q.Where("liquidity_pools.second_coin_id = ?", id)
		} else {
			symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(f.Coin1)
			q = q.Join("LEFT JOIN coins as second_coin").JoinOn("second_coin.id = liquidity_pools.second_coin_id").
				Where(`second_coin.symbol = ?`, symbol).Where(`second_coin.version = ?`, version)
		}
	}

	blocksRange := blocks.RangeSelectFilter{StartBlock: f.StartBlock, EndBlock: f.EndBlock}
	q = q.Apply(blocksRange.Filter)

	return q, nil
}
