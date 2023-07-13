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
