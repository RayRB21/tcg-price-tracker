package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	ID                    string
	Name                  string
	ProductURI            string
	Price                 float64
	PriceChange           int
	PriceChangePercentage float64
	PriceChangeSign       string
	FetchedAt             time.Time
}

type pcProduct struct {
	ID                    string `json:"id"`
	ProductName           string `json:"productName"`
	ProductURI            string `json:"productUri"`
	Price1                string `json:"price1"`
	PriceChange           int    `json:"priceChange"`
	PriceChangePercentage string `json:"priceChangePercentage"`
	PriceChangeSign       string `json:"priceChangeSign"`
}

type pcResponse struct {
	Products []pcProduct `json:"products"`
	Cursor   string      `json:"cursor"`
}

func parsePrice(s string) float64 {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// FetchAllCards paginates through all Riftbound Origins cards
func FetchAllCards() ([]Card, error) {
	var allCards []Card
	cursor := 0
	baseURL := "https://www.pricecharting.com/console/riftbound-origins"

	client := &http.Client{Timeout: 15 * time.Second}

	for {
		url := fmt.Sprintf(
			"%s?sort=&when=none&release-date=2026-03-16&cursor=%d&format=json",
			baseURL, cursor,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Match the header browser sends
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Referer", "https://www.pricecharting.com/")

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch cursor %d: %w", cursor, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status %d at cursor %d", resp.StatusCode, cursor)
		}

		var result pcResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decode cursor %d: %w", cursor, err)
		}

		// No more products
		if len(result.Products) == 0 {
			break
		}

		for _, p := range result.Products {
			pct, _ := strconv.ParseFloat(p.PriceChangePercentage, 64)
			allCards = append(allCards, Card{
				ID:                    p.ID,
				Name:                  p.ProductName,
				ProductURI:            p.ProductURI,
				Price:                 parsePrice(p.Price1),
				PriceChange:           p.PriceChange,
				PriceChangePercentage: pct,
				PriceChangeSign:       p.PriceChangeSign,
				FetchedAt:             time.Now(),
			})
		}

		fmt.Printf("Fetched cursor %d — %d cards so far\n", cursor, len(allCards))

		cursor += 50

		time.Sleep(1 * time.Second)
	}

	return allCards, nil
}
