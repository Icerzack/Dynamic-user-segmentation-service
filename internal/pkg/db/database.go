package db

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate mockgen -source=database.go -destination=mocks/mock.go

type DBops interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	GetPool(ctx context.Context) *pgxpool.Pool
}

type PostgresDatabase struct {
	cluster *pgxpool.Pool
}

func newPostgresDatabase(cluster *pgxpool.Pool) *PostgresDatabase {
	return &PostgresDatabase{cluster: cluster}
}

func (db PostgresDatabase) GetPool(_ context.Context) *pgxpool.Pool {
	return db.cluster
}

func (db PostgresDatabase) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}

func (db PostgresDatabase) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}

func (db PostgresDatabase) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db PostgresDatabase) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

type DBStub struct{}

func NewDBStub() *DBStub {
	return &DBStub{}
}

func (db DBStub) GetPool(_ context.Context) *pgxpool.Pool {
	return nil
}

func (db DBStub) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return nil
}

func (db DBStub) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return nil
}

func (db DBStub) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

func (db DBStub) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return nil
}
