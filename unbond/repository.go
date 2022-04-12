package unbond

import (
	"github.com/MinterTeam/minter-explorer-api/v2/events"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetListByAddress(filter *events.SelectFilter, pagination *tools.Pagination) ([]UnbondMoveStake, error) {
	var unbonds []UnbondMoveStake
	var err error

	allCoins := r.db.Model(new(models.Unbond)).
		ColumnExpr("block_id, coin_id, validator_id, null as to_validator_id, value, created_at, address.address as minter_address").
		Join("JOIN addresses as address ON address.id = unbond.address_id").
		Apply(filter.Filter).
		UnionAll(r.db.Model(new(models.MovedStake)).
			ColumnExpr("block_id, coin_id, from_validator_id as validator_id, to_validator_id, value, created_at, address.address as minter_address").
			Join("JOIN addresses as address ON address.id = moved_stake.address_id").
			Apply(filter.Filter))

	pagination.Total, err = r.db.Model().
		With("data", allCoins).
		Table("data").
		Join("JOIN coins on coins.id = data.coin_id").
		Join("JOIN validators on validators.id = data.validator_id").
		Join("LEFT JOIN validators as toValidators on toValidators.id = data.to_validator_id").
		ColumnExpr("data.block_id, data.coin_id, data.validator_id, data.to_validator_id, data.value, data.created_at").
		ColumnExpr("coins.id as coin__id").
		ColumnExpr("coins.symbol as coin__symbol").
		ColumnExpr("validators.public_key as from_validator__public_key").
		ColumnExpr("validators.name as from_validator__name").
		ColumnExpr("validators.description as from_validator__description").
		ColumnExpr("validators.icon_url as from_validator__icon_url").
		ColumnExpr("validators.site_url as from_validator__site_url").
		ColumnExpr("validators.status as from_validator__status").
		ColumnExpr("toValidators.public_key as to_validator__public_key").
		ColumnExpr("toValidators.name as to_validator__name").
		ColumnExpr("toValidators.description as to_validator__description").
		ColumnExpr("toValidators.icon_url as to_validator__icon_url").
		ColumnExpr("toValidators.site_url as to_validator__site_url").
		ColumnExpr("toValidators.status as to_validator__status").
		ColumnExpr("data.minter_address as address__address").
		OrderExpr("data.block_id desc").
		Apply(pagination.Filter).
		SelectAndCount(&unbonds)

	return unbonds, err
}

func (r *Repository) GetListAsEventsByAddress(filter *events.SelectFilter, lastBlockId uint64, pagination *tools.Pagination) (unbonds []models.Unbond, err error) {
	pagination.Total, err = r.db.Model(&unbonds).
		Relation("Coin").
		Relation("Validator").
		ColumnExpr("unbond.block_id, unbond.value, address.address as address__address").
		Join("JOIN addresses as address ON address.id = unbond.address_id").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Where("unbond.block_id <= ?", lastBlockId).
		Order("unbond.block_id desc").
		SelectAndCount()

	return unbonds, err
}
