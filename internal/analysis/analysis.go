package analysis

import (
	"sort"

	"github.com/RayRB21/tcg-price-tracker/internal/storage"
)

type CardSummary struct {
	CardID      string
	Name        string
	LatestPrice float64
	OldestPrice float64
	ChangeAbs   float64
	ChangePct   float64
	Snapshots   int
}

// Summarise computes price movement for a card across all its snapshots.
func Summarise(name, cardID string, snapshots []storage.PriceSnapshot) CardSummary {
	if len(snapshots) == 0 {
		return CardSummary{CardID: cardID, Name: name}
	}

	oldest := snapshots[0].Price
	latest := snapshots[len(snapshots)-1].Price

	changeAbs := latest - oldest
	changePct := 0.0
	if oldest > 0 {
		changePct = (changeAbs / oldest) * 100
	}

	return CardSummary{
		CardID:      cardID,
		Name:        name,
		LatestPrice: latest,
		OldestPrice: oldest,
		ChangeAbs:   changeAbs,
		ChangePct:   changePct,
		Snapshots:   len(snapshots),
	}
}

// TopMovers returns the N cards with the biggest price change % across runs.
func TopMovers(summaries []CardSummary, n int) []CardSummary {
	// Work on a copy so we don't mutate the original slice
	sorted := make([]CardSummary, len(summaries))
	copy(sorted, summaries)

	sort.Slice(sorted, func(i, j int) bool {
		// Sort by absolute % change so drops show up too
		return abs(sorted[i].ChangePct) > abs(sorted[j].ChangePct)
	})

	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// Spikes returns cards whose latest price jumped by more than threshold %.
func Spikes(summaries []CardSummary, thresholdPct float64) []CardSummary {
	var spikes []CardSummary
	for _, s := range summaries {
		if s.ChangePct >= thresholdPct && s.Snapshots > 1 && s.OldestPrice > 0 {
			spikes = append(spikes, s)
		}
	}

	// Sort biggest spike first
	sort.Slice(spikes, func(i, j int) bool {
		return spikes[i].ChangePct > spikes[j].ChangePct
	})

	return spikes
}

// MovingAverage returns the average price across the last n snapshots.
func MovingAverage(snapshots []storage.PriceSnapshot, n int) float64 {
	if len(snapshots) == 0 {
		return 0
	}
	if n > len(snapshots) {
		n = len(snapshots)
	}

	recent := snapshots[len(snapshots)-n:]
	sum := 0.0
	for _, s := range recent {
		sum += s.Price
	}
	return sum / float64(n)
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
