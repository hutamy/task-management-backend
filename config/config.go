package config

import (
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Port          int    `env:"PORT" envDefault:"8080"`
	DatabaseURL   string `env:"DATABASE_URL"`
	JwtSecret     string `env:"JWT_SECRET"`
	TokenDuration int    `env:"TOKEN_DURATION" envDefault:"24"`
	CacheDuration int    `env:"CACHE_DURATION" envDefault:"24"`
}

var configuration Config

func LoadEnv() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	if err := env.Parse(&configuration); err != nil {
		return Config{}, err
	}

	return configuration, nil
}

func GetConfig() Config {
	return configuration
}

func InitDB(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	tasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		parent_id INTEGER,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (parent_id) REFERENCES tasks(id) ON DELETE CASCADE
	);`

	indexUserID := `CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);`
	indexParentID := `CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id);`

	queries := []string{
		usersTable,
		tasksTable,
		indexUserID,
		indexParentID,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
