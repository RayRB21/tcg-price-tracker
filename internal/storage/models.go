package storage

import "time"

type PriceSnapshot struct {
	Price      float64
	ChangePct  float64
	ChangeSign string
	ScrapedAt  time.Time
}

type CardRecord struct {
	ID   string
	Name string
}
