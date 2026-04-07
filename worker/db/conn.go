package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func Connect(ctx context.Context, connection string) error {
	c, err := pgx.Connect(ctx, connection)
	if err != nil {
		return fmt.Errorf("unable to open database connection: %w", err)
	}

	conn = c

	return nil
}

func Close(ctx context.Context) {
	defer conn.Close(ctx)
}

func Connection() *pgx.Conn {
	return conn
}
