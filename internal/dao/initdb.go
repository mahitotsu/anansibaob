package dao

import (
	"context"
	"fmt"
)

type dbInitializer struct {
	SqlClient SqlClient
}

type DBInitializer interface {
	DropAndCreate(ctx context.Context) error
}

var _ DBInitializer = (*dbInitializer)(nil)

func NewDBInitializer(sqlClient SqlClient) DBInitializer {
	return &dbInitializer{SqlClient: sqlClient}
}

func (i *dbInitializer) DropAndCreate(ctx context.Context) error {

	if tableNames, err := i.GetAllTableNames(ctx); err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	} else {
		for _, tableName := range tableNames {
			if _, err := i.DropTable(ctx, tableName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *dbInitializer) DropTable(ctx context.Context, tableName string) (int64, error) {

	if tx, err := i.SqlClient.Begin(ctx); err != nil {
		return 0, err
	} else {
		defer tx.Rollback(ctx)
		if num, err := tx.Exec(ctx, "DROP TABLE IF EXISTS "+tableName); err != nil {
			return 0, err
		} else {
			tx.Commit(ctx)
			return num, nil
		}
	}
}

func (i *dbInitializer) GetAllTableNames(ctx context.Context) ([]string, error) {

	var tableNames []string
	if tx, err := i.SqlClient.Begin(ctx); err != nil {
		return nil, err
	} else {
		defer tx.Rollback(ctx)
		if rows, err := tx.Query(ctx, `
			SELECT tablename FROM pg_tables
			WHERE schemaname = 'public'
			AND schemaname NOT LIKE 'pg_%'
			AND schemaname NOT LIKE 'sql_%'
		`); err != nil {
			return nil, err
		} else {
			for _, row := range rows {
				tableNames = append(tableNames, row["tablename"].(string))
			}
		}
	}
	return tableNames, nil
}
