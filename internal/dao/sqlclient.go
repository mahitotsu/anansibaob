package dao

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sqlclient struct {
	pool *pgxpool.Pool
}

type SqlClient interface {
	Query(ctx context.Context, statement string) ([]map[string]any, error)
}

var _ SqlClient = (*sqlclient)(nil)

func NewSqlClient(pool *pgxpool.Pool) SqlClient {
	return &sqlclient{pool: pool}
}

func (c *sqlclient) Query(ctx context.Context, statement string) ([]map[string]any, error) {

	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToMap)
}
