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
	Type Type
}

func NewTypeFilter(f string) TypeFilter {
	return TypeFilter{}
}

func (f TypeFilter) Filter(q *orm.Query) (*orm.Query, error) {
	return q.Where(""), nil
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
