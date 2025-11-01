package main

import (
	"context"
	"database/sql"
	"flag"
	"log"

	migrations "github.com/silent-observer/go-tickets/pgdb/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
)

var postgres_conn = flag.String("postgres_conn", 
		"postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
		"PostgreSQL connection string")
var rollback = flag.Bool("rollback", false, "Perform rollback")

func main() {
	cxt := context.Background()

	flag.Parse()
	log.Println("Starting...")
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(*postgres_conn)))
	db := bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	migrator := migrate.NewMigrator(db, migrations.Migrations)
	migrator.Init(cxt)
	err := migrator.Lock(cxt)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	defer migrator.Unlock(cxt)

	if (*rollback) {
		migrator.Rollback(cxt)
	} else {
		migrator.Migrate(cxt)
	}
}