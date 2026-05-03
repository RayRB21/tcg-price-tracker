package main

import (
	"fmt"
	"log"

	"github.com/RayRB21/tcg-price-tracker/internal/analysis"
	"github.com/RayRB21/tcg-price-tracker/internal/scraper"
	"github.com/RayRB21/tcg-price-tracker/internal/storage"
)

func main() {
	db, err := storage.Open("prices.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Fetching cards...")
	cards, err := scraper.FetchAllCards()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d cards — saving snapshots...\n\n", len(cards))
	for _, c := range cards {
		if err := db.SaveSnapshot(c); err != nil {
			log.Printf("warning: %v", err)
		}
	}

	// --- Analysis ---
	fmt.Println("\n=== Loading price history for analysis ===\n")

	allCards, err := db.GetAllCards()
	if err != nil {
		log.Fatal(err)
	}

	var summaries []analysis.CardSummary
	for _, c := range allCards {
		snapshots, err := db.GetSnapshots(c.ID)
		if err != nil {
			log.Printf("warning: %v", err)
			continue
		}
		summary := analysis.Summarise(c.Name, c.ID, snapshots)
		summaries = append(summaries, summary)
	}

	// Top 10 movers across all your runs
	fmt.Println("=== TOP 10 MOVERS (across all runs) ===")
	for i, s := range analysis.TopMovers(summaries, 10) {
		sign := "+"
		if s.ChangePct < 0 {
			sign = ""
		}
		fmt.Printf("%2d. %-50s $%.2f → $%.2f  (%s%.1f%%)\n",
			i+1, s.Name, s.OldestPrice, s.LatestPrice, sign, s.ChangePct)
	}

	// Spikes — cards up more than 10% across your runs
	fmt.Println("\n=== SPIKES (>10% increase across runs) ===")
	spikes := analysis.Spikes(summaries, 10.0)
	if len(spikes) == 0 {
		fmt.Println("No spikes detected — run the scraper a few more times to build history.")
	} else {
		for _, s := range spikes {
			fmt.Printf("  %-50s $%.2f → $%.2f  (+%.1f%%)\n",
				s.Name, s.OldestPrice, s.LatestPrice, s.ChangePct)
		}
	}

	// Moving average for a specific card — Wallop Foil since it was spiking
	fmt.Println("\n=== MOVING AVERAGE (Wallop Foil) ===")
	for _, c := range allCards {
		if c.Name == "Wallop [Foil] #146" {
			snapshots, _ := db.GetSnapshots(c.ID)
			avg := analysis.MovingAverage(snapshots, 3)
			latest := snapshots[len(snapshots)-1].Price
			fmt.Printf("  Last 3-run average: $%.2f  |  Latest: $%.2f\n", avg, latest)
			fmt.Printf("  Based on %d snapshots\n", len(snapshots))
		}
	}
}
