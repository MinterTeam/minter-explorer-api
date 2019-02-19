package tools

import "github.com/go-pg/pg/orm"

type Filter interface {
	Filter(q *orm.Query) (*orm.Query, error)
}
