package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := migrate(conn); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func migrate(conn *sql.DB) error {
	_, err := conn.Exec(`
        CREATE TABLE IF NOT EXISTS cards (
            id          TEXT PRIMARY KEY,
            name        TEXT NOT NULL,
            product_uri TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS price_snapshots (
            id          INTEGER PRIMARY KEY AUTOINCREMENT,
            card_id     TEXT NOT NULL REFERENCES cards(id),
            price       REAL NOT NULL,
            change_pct  REAL NOT NULL,
            change_sign TEXT NOT NULL,
            scraped_at  DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE INDEX IF NOT EXISTS idx_snapshots_card_time
            ON price_snapshots(card_id, scraped_at);
    `)
	return err
}
