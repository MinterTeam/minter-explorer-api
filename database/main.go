package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/go-pg/pg/v9"
	"os"
)

func Connect(env *core.Environment) *pg.DB {
	options := &pg.Options{
		User:     env.DbUser,
		Password: env.DbPassword,
		Database: env.DbName,
		PoolSize: env.DbPoolSize,
		Addr:     fmt.Sprintf("%s:%s", env.DbHost, env.DbPort),
	}

	if os.Getenv("POSTGRES_SSL_ENABLED") == "true" {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	db := pg.Connect(options)
	if db == nil {
		panic("Could not connect to database")
	}

	if env.IsDebug {
		db.AddQueryHook(dbLogger{})
	}

	return db
}

func Close(db *pg.DB) {
	err := db.Close()
	if err != nil {
		panic(fmt.Sprintf("Could not close connection to database: %s", err))
	}
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}
