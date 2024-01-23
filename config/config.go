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
	// schema, err := os.ReadFile("models/schema/schema.sql")
	// if err != nil {
	// 	log.Fatalf("Error reading schema.sql: %v", err)
	// }
	// _, err = conn.Exec(context.Background(), string(schema))
	// if err != nil {
	// 	log.Fatalf("Error executing schema.sql: %v", err)
	// }
	// defer conn.Close(ctx)
	return conn
}
