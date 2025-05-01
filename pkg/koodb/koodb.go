package koodb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// NewPostgresPool creates a new sql.DB handle that uses pgxpool
// which is a connection pool for PostgreSQL.
func NewPostgresPool(ctx context.Context, dsn string) (*sql.DB, error) {
	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set default query exec mode to simple protocol as bun does not benefit from the prepared statement cache
	pgxPoolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	sqldb := stdlib.OpenDBFromPool(pgxPool)

	return sqldb, nil
}

// NewPostgresTxDB creates a new sql.DB handle that uses txdb
// which wraps the entire connection in a transaction.
// This is useful for testing purposes where we want to
// ensure that the database is always in a known state.
func NewPostgresTxDB(dsn string) (*sql.DB, error) {
	txDB := sql.OpenDB(txdb.New("pgx", dsn))

	return txDB, nil
}

// NewPostgresConn creates a new sql.DB handle that uses pgx
// which is a single connection for PostgreSQL.
func NewPostgresConn(ctx context.Context, dsn string) (*sql.DB, error) {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqldb := stdlib.OpenDB(*config)

	return sqldb, nil
}
