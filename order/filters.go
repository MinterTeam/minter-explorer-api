package order

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10/orm"
)

type AddressFilter struct {
	Address string
}

func NewAddressFilter(address string) AddressFilter {
	return AddressFilter{address}
}

func (f AddressFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if len(f.Address) == 0 {
		return q, nil
	}

	return q.Where(`"address"."address" = ?`, f.Address), nil
}

// ------------------------------

type TypeFilter struct {
	Type   Type
	pool   models.LiquidityPool
	coinId uint64
}

func NewTypeFilter(f string, pool models.LiquidityPool, coinId uint64) TypeFilter {
	return TypeFilter{Type(f), pool, coinId}
}

func (f TypeFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.Type == OrderTypeBuy {
		q = q.Where("coin_buy_id = ?", f.coinId)

		if f.pool.FirstCoinId == f.coinId {
			q = q.OrderExpr("price desc")
		} else {
			q = q.OrderExpr("price asc")
		}
	}

	if f.Type == OrderTypeSell {
		q = q.Where("coin_sell_id = ?", f.coinId)

		if f.pool.FirstCoinId == f.coinId {
			q = q.OrderExpr("price asc")
		} else {
			q = q.OrderExpr("price desc")
		}
	}

	if len(f.Type) == 0 {
		return q.OrderExpr(`"id" desc`), nil
	}

	return q, nil
}

// ------------------------------

type PoolFilter struct {
	Pool models.LiquidityPool
}

func NewPoolFilter(p models.LiquidityPool) PoolFilter {
	return PoolFilter{p}
}

func (f PoolFilter) Filter(q *orm.Query) (*orm.Query, error) {
	return q.Where("liquidity_pool_id = ?", f.Pool.Id), nil
}


// ------------------------------

type StatusFilter struct {
	status Status
}

func NewStatusFilter(status string) StatusFilter {
	return StatusFilter{Status(status)}
}

func (f StatusFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if len(f.status) == 0 {
		return q, nil
	}

	if f.status == OrderStatusActive {
		return q.Where(`"status" in (?, ?)`, models.OrderTypeActive, models.OrderTypeNew), nil
	}

	if f.status == OrderStatusCanceled {
		return q.Where(`"status" = ?`, models.OrderTypeCanceled), nil
	}

	if f.status == OrderStatusExpired {
		return q.Where(`"status" = ?`, models.OrderTypeExpired), nil
	}

	if f.status == OrderStatusFilled {
		return q.Where(`"status" = ?`, models.OrderTypeFilled), nil
	}

	if f.status == OrderStatusPartiallyFilled {
		return q.Where(`"status" = ?`, models.OrderTypePartiallyFilled), nil
	}

	return q, nil
}
