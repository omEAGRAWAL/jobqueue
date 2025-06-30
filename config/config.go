package config

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitDB() *sqlx.DB {
	dsn := "postgresql://neondb_owner:npg_VzoF2mdZ0hXp@ep-green-morning-a82bpjj4.eastus2.azure.neon.tech/neondb?sslmode=require"
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DATABASE_URL: %v", err)
	}

	cfg.RuntimeParams["sslmode"] = "require"
	cfg.RuntimeParams["application_name"] = "jobqueue-app"

	db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	return db
}
