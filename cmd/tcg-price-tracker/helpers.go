package main

import (
	"log"

	"github.com/RayRB21/tcg-price-tracker/internal/analysis"
	"github.com/RayRB21/tcg-price-tracker/internal/storage"
)

func buildSummaries(db *storage.DB) ([]analysis.CardSummary, error) {
	allCards, err := db.GetAllCards()
	if err != nil {
		return nil, err
	}

	var summaries []analysis.CardSummary
	for _, c := range allCards {
		snapshots, err := db.GetSnapshots(c.ID)
		if err != nil {
			log.Printf("warning: skipping %s: %v", c.Name, err)
			continue
		}
		summaries = append(summaries, analysis.Summarise(c.Name, c.ID, snapshots))
	}
	return summaries, nil
}
