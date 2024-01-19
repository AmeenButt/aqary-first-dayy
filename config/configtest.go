package config

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func TestConfig(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, "postgresql://postgres:admin@localhost:5432/test_sqlc_practice")
	if err != nil {
		panic(err)
	}
	// defer conn.Close(ctx)
	return conn
}
