package storage

import (
	"fmt"

	"github.com/RayRB21/tcg-price-tracker/internal/scraper"
)

// SaveSnapshot inserts or ignores the card record, then saves a price snapshot.
func (db *DB) SaveSnapshot(c scraper.Card) error {
	_, err := db.conn.Exec(`
        INSERT OR IGNORE INTO cards (id, name, product_uri)
        VALUES (?, ?, ?)`,
		c.ID, c.Name, c.ProductURI,
	)
	if err != nil {
		return fmt.Errorf("upsert card %s: %w", c.Name, err)
	}

	// Don't save duplicate snapshots within the same hour
	var count int
	err = db.conn.QueryRow(`
        SELECT COUNT(*) FROM price_snapshots
        WHERE card_id = ?
        AND scraped_at > datetime('now', '-1 hour')`,
		c.ID,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("check duplicate %s: %w", c.Name, err)
	}
	if count > 0 {
		return nil // already saved this card recently, skip
	}

	_, err = db.conn.Exec(`
        INSERT INTO price_snapshots (card_id, price, change_pct, change_sign)
        VALUES (?, ?, ?, ?)`,
		c.ID, c.Price, c.PriceChangePercentage, c.PriceChangeSign,
	)
	if err != nil {
		return fmt.Errorf("insert snapshot %s: %w", c.Name, err)
	}

	return nil
}

// GetSnapshots returns all price snapshots for a card, oldest first.
func (db *DB) GetSnapshots(cardID string) ([]PriceSnapshot, error) {
	rows, err := db.conn.Query(`
        SELECT price, change_pct, change_sign, scraped_at
        FROM price_snapshots
        WHERE card_id = ?
        ORDER BY scraped_at ASC`,
		cardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []PriceSnapshot
	for rows.Next() {
		var s PriceSnapshot
		if err := rows.Scan(&s.Price, &s.ChangePct, &s.ChangeSign, &s.ScrapedAt); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (db *DB) GetAllCards() ([]CardRecord, error) {
	rows, err := db.conn.Query(`SELECT id, name FROM cards`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CardRecord
	for rows.Next() {
		var c CardRecord
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}
