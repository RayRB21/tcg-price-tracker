package main

import (
	"fmt"
	"log"

	"github.com/RayRB21/tcg-price-tracker/internal/scraper"
	"github.com/RayRB21/tcg-price-tracker/internal/storage"
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Fetch latest Riftbound prices and save to database",
	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Printf("Fetched %d cards — saving snapshots...\n", len(cards))
		saved := 0
		for _, c := range cards {
			if err := db.SaveSnapshot(c); err != nil {
				log.Printf("warning: %v", err)
				continue
			}
			saved++
		}
		fmt.Printf("Done — saved %d snapshots.\n", saved)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
