package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn
var queries *Queries

func Connect(ctx context.Context, connection string) error {
	c, err := pgx.Connect(ctx, connection)
	if err != nil {
		return fmt.Errorf("unable to open database connection: %w", err)
	}

	conn = c

	queries = New(conn)

	return nil
}

func Close(ctx context.Context) {
	defer conn.Close(ctx)
}

var mutex sync.Mutex

func Get() (*Queries, *pgx.Conn) {
	mutex.Lock()

	return queries, conn
}

func Release() {
	mutex.Unlock()
}
