package configs

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// NewMySQLConnection creates a new database connection
func NewMySQLConnection(env *Env, zerolog *zerolog.Logger) *sqlx.DB {
	_ = zerolog.With().Str("method", "MySQLConnection").Logger()

	dsn := getMySQLConnectionString(env)

	// Setup connection
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to MySQL database: %w", err))
	}

	if err := configureConnection(db); err != nil {
		db.Close()
		panic(fmt.Errorf("failed to connect to MySQL database: %w", err))
	}

	runMigrations(dsn)
	log.Println("Connected to MySQL")
	return db
}

// configureConnection sets up the connection pool settings for both connections
func configureConnection(dbs ...*sqlx.DB) error {
	for _, db := range dbs {
		db.SetMaxIdleConns(20)
		db.SetMaxOpenConns(200)
		db.SetConnMaxLifetime(time.Hour)

		if err := db.Ping(); err != nil {
			return fmt.Errorf("database ping failed: %w", err)
		}
	}
	return nil
}

// getMySQLConnectionString is used to get the MySQL connection string
func getMySQLConnectionString(env *Env) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s&charset=utf8mb4&collation=utf8mb4_unicode_ci&multiStatements=true",
		env.Get("MYSQL_USER"),
		env.Get("MYSQL_PASSWORD"),
		env.Get("MYSQL_HOST"),
		env.Get("MYSQL_PORT"),
		env.Get("MYSQL_DATABASE"),
		url.QueryEscape(env.Get("TIMEZONE")),
	)
}
