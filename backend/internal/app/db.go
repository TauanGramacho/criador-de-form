package app

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
	dialect string
}

func OpenDB(databaseURL string) (*Database, error) {
	dialect := detectDialect(databaseURL)
	if dialect == "sqlite" {
		if err := ensureSQLiteDir(databaseURL); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open(driverName(dialect), databaseURL)
	if err != nil {
		return nil, err
	}
	if dialect == "sqlite" {
		db.SetMaxOpenConns(1)
		if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
			_ = db.Close()
			return nil, err
		}
	} else {
		db.SetMaxOpenConns(5)
		db.SetMaxIdleConns(5)
	}
	return &Database{DB: db, dialect: dialect}, nil
}

func detectDialect(databaseURL string) string {
	if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		return "postgres"
	}
	return "sqlite"
}

func driverName(dialect string) string {
	if dialect == "postgres" {
		return "pgx"
	}
	return "sqlite"
}

func (db *Database) Rebind(query string) string {
	if db.dialect != "postgres" {
		return query
	}
	var builder strings.Builder
	argIndex := 1
	for _, char := range query {
		if char == '?' {
			builder.WriteString(fmt.Sprintf("$%d", argIndex))
			argIndex++
			continue
		}
		builder.WriteRune(char)
	}
	return builder.String()
}

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	return db.DB.Exec(db.Rebind(query), args...)
}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	return db.DB.Query(db.Rebind(query), args...)
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	return db.DB.QueryRow(db.Rebind(query), args...)
}

func ensureSQLiteDir(databaseURL string) error {
	if !strings.HasPrefix(databaseURL, "file:") {
		return nil
	}
	path := strings.TrimPrefix(databaseURL, "file:")
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}
	if path == "" || path == ":memory:" {
		return nil
	}
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func RunMigrations(db *Database) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		return err
	}

	migrationsPath, err := findMigrationsPath()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		var existing string
		err := db.QueryRow(`SELECT version FROM schema_migrations WHERE version = ?`, name).Scan(&existing)
		if err == nil {
			continue
		}
		if err != sql.ErrNoRows {
			return err
		}

		sqlBytes, err := os.ReadFile(filepath.Join(migrationsPath, name))
		if err != nil {
			return err
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply %s: %w", name, err)
		}
		if _, err := tx.Exec(db.Rebind(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)`), name, nowUTC().Format(timeFormat)); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func findMigrationsPath() (string, error) {
	candidates := []string{"migrations", "backend/migrations"}
	for _, candidate := range candidates {
		if entries, err := os.ReadDir(candidate); err == nil && len(entries) > 0 {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("migrations directory not found")
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
