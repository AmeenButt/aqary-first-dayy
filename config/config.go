package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Config(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	// defer conn.Close(ctx)
	return conn
}
