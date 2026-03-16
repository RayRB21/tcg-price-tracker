package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/RayRB21/tcg-price-tracker/internal/scraper"
)

func main() {
	cards, err := scraper.FetchAllCards()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nFetched %d cards total\n\n", len(cards))
	for _, c := range cards {
		fmt.Printf("%-50s $%.2f  %s%s%%\n",
			c.Name, c.Price, c.PriceChangeSign,
			strconv.FormatFloat(c.PriceChangePercentage, 'f', 1, 64))
	}
}
