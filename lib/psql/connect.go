package psql

import (
	"context"
	"fmt"
	"time"

	"cp/lib/cfg"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitPgxPool(c *cfg.Config) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.SqlUser,
		c.SqlPass,
		c.SqlHost,
		c.SqlPort,
		c.SqlDb))

	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
