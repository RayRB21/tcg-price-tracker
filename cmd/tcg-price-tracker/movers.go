package main

import (
	"fmt"
	"log"

	"github.com/RayRB21/tcg-price-tracker/internal/analysis"
	"github.com/RayRB21/tcg-price-tracker/internal/storage"
	"github.com/spf13/cobra"
)

var topN int

var moversCmd = &cobra.Command{
	Use:   "movers",
	Short: "Show the top N cards by price movement across runs",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := storage.Open("prices.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		summaries, err := buildSummaries(db)
		if err != nil {
			log.Fatal(err)
		}

		movers := analysis.TopMovers(summaries, topN)
		fmt.Printf("=== TOP %d MOVERS ===\n\n", topN)
		for i, s := range movers {
			sign := "+"
			if s.ChangePct < 0 {
				sign = ""
			}
			fmt.Printf("%2d. %-50s $%.2f → $%.2f  (%s%.1f%%)\n",
				i+1, s.Name, s.OldestPrice, s.LatestPrice, sign, s.ChangePct)
		}
	},
}

func init() {
	moversCmd.Flags().IntVarP(&topN, "top", "n", 10, "Number of cards to show")
	rootCmd.AddCommand(moversCmd)
}
