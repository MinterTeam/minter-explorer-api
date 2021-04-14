package tools

import "github.com/go-pg/pg/v10/orm"

type Filter interface {
	Filter(q *orm.Query) (*orm.Query, error)
}
