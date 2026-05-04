package main

import (
	"fmt"
	"log"

	"github.com/RayRB21/tcg-price-tracker/internal/analysis"
	"github.com/RayRB21/tcg-price-tracker/internal/storage"
	"github.com/spf13/cobra"
)

var cardName string

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show price history for a specific card",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := storage.Open("prices.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		allCards, err := db.GetAllCards()
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range allCards {
			if c.Name != cardName {
				continue
			}

			snapshots, err := db.GetSnapshots(c.ID)
			if err != nil {
				log.Fatal(err)
			}

			avg := analysis.MovingAverage(snapshots, 3)
			fmt.Printf("=== %s ===\n\n", c.Name)
			fmt.Printf("  Snapshots:       %d\n", len(snapshots))
			fmt.Printf("  Oldest price:    $%.2f\n", snapshots[0].Price)
			fmt.Printf("  Latest price:    $%.2f\n", snapshots[len(snapshots)-1].Price)
			fmt.Printf("  3-run avg:       $%.2f\n\n", avg)
			fmt.Println("  Full history:")
			for _, s := range snapshots {
				fmt.Printf("    %s  $%.2f\n", s.ScrapedAt.Format("2006-01-02 15:04"), s.Price)
			}
			return
		}

		fmt.Printf("Card not found: %q\n", cardName)
		fmt.Println("Tip: use the exact name from the movers or spikes output.")
	},
}

func init() {
	historyCmd.Flags().StringVarP(&cardName, "card", "c", "", "Card name to look up (required)")
	historyCmd.MarkFlagRequired("card")
	rootCmd.AddCommand(historyCmd)
}
