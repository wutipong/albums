package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool
var queries *Queries

func Connect(ctx context.Context, connection string) error {
	conn, err := pgxpool.New(ctx, connection)
	if err != nil {
		return fmt.Errorf("unable to open database connection: %w", err)
	}

	queries = New(conn)

	return nil
}

func Close(ctx context.Context) {
	defer conn.Close()
}

func Get() (*Queries, *pgxpool.Pool) {
	return queries, conn
}
