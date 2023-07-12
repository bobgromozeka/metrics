package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var pgxConnection *pgx.Conn

func Connect(dsn string) error {
	if conn, connErr := pgx.Connect(context.Background(), dsn); connErr != nil {
		return connErr
	} else {
		pgxConnection = conn
	}

	return nil
}

func Connection() *pgx.Conn {
	return pgxConnection
}

func ExecTablesDDL() error {
	ctx := context.Background()
	tx, txErr := pgxConnection.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	_, gaugeErr := tx.Exec(
		context.Background(),
		`create table if not exists gauges(
    			name varchar(255) not null,
    			value double precision,
    			primary key (name)
    			)`,
	)
	if gaugeErr != nil {
		tx.Rollback(ctx)
		return gaugeErr
	}

	_, counterErr := tx.Exec(
		context.Background(),
		`create table if not exists counters(
    			name varchar(255) not null,
    			value bigint,
    			primary key (name)
    			)`,
	)
	if counterErr != nil {
		tx.Rollback(ctx)
		return counterErr
	}

	tx.Commit(ctx)

	return nil
}
