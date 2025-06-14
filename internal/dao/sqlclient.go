package dao

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sqlclient struct {
	pool *pgxpool.Pool
}

type sqltx struct {
	tx pgx.Tx
}

type SqlClient interface {
	Begin(ctx context.Context) (SqlTx, error)
}

type SqlTx interface {
	Query(ctx context.Context, statement string) ([]map[string]any, error)
	Exec(ctx context.Context, statement string) (int64, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

var _ SqlClient = (*sqlclient)(nil)
var _ SqlTx = (*sqltx)(nil)

func NewSqlClient(pool *pgxpool.Pool) SqlClient {
	return &sqlclient{pool: pool}
}

func (c *sqlclient) Begin(ctx context.Context) (SqlTx, error) {
	if tx, err := c.pool.Begin(ctx); err != nil {
		return nil, err
	} else {
		return &sqltx{tx: tx}, nil
	}
}

func (c *sqltx) Query(ctx context.Context, statement string) ([]map[string]any, error) {

	if rows, err := c.tx.Query(ctx, statement); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		return pgx.CollectRows(rows, pgx.RowToMap)
	}
}

func (c *sqltx) Exec(ctx context.Context, statement string) (int64, error) {

	if tag, err := c.tx.Exec(ctx, statement); err != nil {
		return 0, err
	} else {
		return tag.RowsAffected(), nil
	}
}

func (c *sqltx) Commit(ctx context.Context) error {
	return c.tx.Commit(ctx)
}

func (c *sqltx) Rollback(ctx context.Context) error {
	return c.tx.Rollback(ctx)
}
