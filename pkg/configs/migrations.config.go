package configs

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dsn string) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Could not create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	// Check if the migration is in a dirty state and force if necessary
	_, _, err = m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Could not get migration version: %v", err)
	}

	// !NOTE: this can delete all your data from the database so use with caution
	// if dirty {
	// 	if err := m.Force(int(version)); err != nil {
	// 		log.Fatalf("Could not force migration version: %v", err)
	// 	}
	// }

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %v", err)
	}
	log.Println("Migrations ran successfully")
}
