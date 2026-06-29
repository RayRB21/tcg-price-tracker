package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/RayRB21/tcg-price-tracker/internal/storage"
	"github.com/spf13/cobra"
)

var exportPath string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export full price history to a CSV file",
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

		f, err := os.Create(exportPath)
		if err != nil {
			log.Fatalf("create file: %v", err)
		}
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()

		// Header row
		w.Write([]string{"Card Name", "Price (USD)", "Change %", "Change Direction", "Scraped At"})

		for _, c := range allCards {
			snapshots, err := db.GetSnapshots(c.ID)
			if err != nil {
				log.Printf("warning: skipping %s: %v", c.Name, err)
				continue
			}
			for _, s := range snapshots {
				w.Write([]string{
					c.Name,
					strconv.FormatFloat(s.Price, 'f', 2, 64),
					strconv.FormatFloat(s.ChangePct, 'f', 2, 64),
					s.ChangeSign,
					s.ScrapedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}

		fmt.Printf("Exported to %s\n", exportPath)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportPath, "output", "o", "riftbound_prices.csv", "Output file path")
	rootCmd.AddCommand(exportCmd)
}
