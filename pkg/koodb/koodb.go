package koodb

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// PoolConfig holds configuration for the database connection pool.
type PoolConfig struct {
	MaxConns          int           // Maximum number of connections in the pool (default: 25)
	MinConns          int           // Minimum number of connections in the pool (default: 5)
	MaxConnLifetime   time.Duration // Maximum lifetime of a connection (default: 1 hour)
	MaxConnIdleTime   time.Duration // Maximum idle time of a connection (default: 30 minutes)
	ConnectionTimeout time.Duration // Timeout for establishing a connection (default: 10 seconds)
}

// NewPostgresPool creates a new sql.DB handle with custom pool configuration.
func NewPostgresPool(ctx context.Context, dsn string, poolConfig *PoolConfig) (*sql.DB, error) {
	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set default query exec mode to simple protocol as bun does not benefit from the prepared statement cache
	pgxPoolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Apply pool configuration if provided
	if poolConfig != nil {
		if poolConfig.MaxConns > 0 && poolConfig.MaxConns <= math.MaxInt32 {
			pgxPoolConfig.MaxConns = int32(poolConfig.MaxConns)
		}

		if poolConfig.MinConns > 0 && poolConfig.MinConns <= math.MaxInt32 {
			pgxPoolConfig.MinConns = int32(poolConfig.MinConns)
		}

		if poolConfig.MaxConnLifetime > 0 {
			pgxPoolConfig.MaxConnLifetime = poolConfig.MaxConnLifetime
		}

		if poolConfig.MaxConnIdleTime > 0 {
			pgxPoolConfig.MaxConnIdleTime = poolConfig.MaxConnIdleTime
		}

		if poolConfig.ConnectionTimeout > 0 {
			pgxPoolConfig.ConnConfig.ConnectTimeout = poolConfig.ConnectionTimeout
		}
	}

	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Verify connectivity
	if err := pgxPool.Ping(ctx); err != nil {
		pgxPool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
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
