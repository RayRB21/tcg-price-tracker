package main

import (
	"fmt"
	"log"

	"github.com/RayRB21/tcg-price-tracker/internal/analysis"
	"github.com/RayRB21/tcg-price-tracker/internal/storage"
	"github.com/spf13/cobra"
)

var threshold float64

var spikesCmd = &cobra.Command{
	Use:   "spikes",
	Short: "Show cards that have spiked in price across runs",
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

		spikes := analysis.Spikes(summaries, threshold)
		if len(spikes) == 0 {
			fmt.Printf("No spikes detected above %.1f%%\n", threshold)
			return
		}

		fmt.Printf("=== SPIKES (>%.0f%% increase across runs) ===\n\n", threshold)
		for _, s := range spikes {
			fmt.Printf("  %-50s $%.2f → $%.2f  (+%.1f%%)\n",
				s.Name, s.OldestPrice, s.LatestPrice, s.ChangePct)
		}
	},
}

func init() {
	spikesCmd.Flags().Float64VarP(&threshold, "threshold", "t", 10.0,
		"Minimum % increase to count as a spike")
	rootCmd.AddCommand(spikesCmd)
}
