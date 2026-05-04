# TCG Price Tracker

A CLI tool written in Go that tracks **Riftbound** trading card prices over time, detecting price spikes and trends across the full card catalogue.

Prices are sourced by reverse engineering PriceCharting's internal JSON endpoint — discovered via Chrome DevTools network analysis — and stored locally in a SQLite database. Each run adds a new price snapshot per card, building up a time-series history that the analysis engine queries to surface meaningful market movements.

---

## Features

- **Full catalogue scraping** — fetches all 700+ Riftbound Origins cards with live prices via paginated JSON endpoint
- **Time-series storage** — every scrape saves a timestamped snapshot to SQLite, building price history over time
- **Spike detection** — surfaces cards whose price has increased beyond a configurable threshold across runs
- **Top movers** — ranks cards by absolute price movement (up or down) across all stored snapshots
- **Price history** — shows the full price timeline for any individual card with moving average
- **Deduplication** — prevents duplicate snapshots within the same hour

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.21+ |
| CLI framework | [cobra](https://github.com/spf13/cobra) |
| HTML/JSON fetching | `net/http` + `encoding/json` |
| Database | SQLite via [go-sqlite3](https://github.com/mattn/go-sqlite3) |
| Data source | PriceCharting.com (reverse-engineered internal JSON endpoint) |

---

## Project Structure

```
tcg-price-tracker/
├── cmd/tcg-price-tracker/
│   ├── main.go          # Entry point
│   ├── root.go          # Cobra root command
│   ├── scrape.go        # scrape subcommand
│   ├── movers.go        # movers subcommand
│   ├── spikes.go        # spikes subcommand
│   ├── history.go       # history subcommand
│   └── helpers.go       # Shared analysis helpers
├── internal/
│   ├── scraper/
│   │   └── scraper.go   # HTTP client + JSON pagination
│   ├── storage/
│   │   ├── db.go        # SQLite connection + migrations
│   │   ├── queries.go   # SaveSnapshot, GetSnapshots, GetAllCards
│   │   └── models.go    # Shared types
│   └── analysis/
│       └── analysis.go  # Spike detection, top movers, moving average
├── Makefile
├── go.mod
└── go.sum
```

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- GCC (required for go-sqlite3 CGo compilation)
  - **Windows:** Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
  - **Mac:** `xcode-select --install`
  - **Linux:** `sudo apt install gcc`

---

## Setup

```bash
# Clone the repository
git clone https://github.com/RayRB21/tcg-price-tracker
cd tcg-price-tracker

# Install dependencies
go mod tidy

# Build the binary
make build
# or: go build -o bin/tcg ./cmd/tcg-price-tracker
```

---

## Usage

### Scrape latest prices

Fetches all cards from PriceCharting and saves a snapshot to the local database.

```bash
make scrape
# or: go run ./cmd/tcg-price-tracker scrape
```

```
Fetching cards...
Fetched cursor 0 — 150 cards so far
Fetched cursor 50 — 300 cards so far
...
Fetched 2050 cards — saving snapshots...
Done — saved 2050 snapshots.
```

---

### Top movers

Shows the N cards with the biggest price movement across all stored runs.

```bash
go run ./cmd/tcg-price-tracker movers --top 10
```

```
=== TOP 10 MOVERS ===

 1. Poro Herder [Foil] #61                             $0.15 → $0.41  (+173.3%)
 2. Sett - The Boss #269                               $0.31 → $0.66  (+112.9%)
 3. Confront #129                                      $0.19 → $0.40  (+110.5%)
 4. Sprite Call [Foil] #94                             $0.94 → $1.82  (+93.6%)
 5. Stealthy Pursuer [Promo] #177                      $0.09 → $0.17  (+88.9%)
 6. Sky Splitter #14                                   $0.11 → $0.20  (+81.8%)
 7. Blazing Scorcher #1                                $0.21 → $0.08  (-61.9%)
 8. Bandle Tree [Foil] #278                            $0.20 → $0.32  (+60.0%)
 9. Malzahar - Fanatic #113                            $0.46 → $0.73  (+58.7%)
10. Champion Deck: Jinx                                $15.63 → $24.63  (+57.6%)
```

---

### Spike detection

Surfaces cards that have increased beyond a price threshold across runs.

```bash
go run ./cmd/tcg-price-tracker spikes --threshold 20
```

```
=== SPIKES (>20% increase across runs) ===

  Poro Herder [Foil] #61                             $0.15 → $0.41  (+173.3%)
  Sett - The Boss #269                               $0.31 → $0.66  (+112.9%)
  Confront #129                                      $0.19 → $0.40  (+110.5%)
  ...
```

---

### Price history

Shows the full price timeline for a specific card, including a 3-run moving average.

```bash
go run ./cmd/tcg-price-tracker history --card "Wallop [Foil] #146"
```

```
=== Wallop [Foil] #146 ===

  Snapshots:       9
  Oldest price:    $0.67
  Latest price:    $1.64
  3-run avg:       $1.55

  Full history:
    2026-04-29 17:48  $0.67
    2026-04-30 09:12  $1.12
    2026-05-01 00:15  $1.64
    ...
```

---

### All commands

```bash
go run ./cmd/tcg-price-tracker --help
```

```
Tracks Riftbound card prices over time and detects trends and spikes.

Usage:
  tcg [command]

Available Commands:
  scrape      Fetch latest Riftbound prices and save to database
  movers      Show the top N cards by price movement across runs
  spikes      Show cards that have spiked in price across runs
  history     Show price history for a specific card
  help        Help about any command

Flags:
  -h, --help   help for tcg
```

---

## How It Works

1. **Endpoint discovery** — PriceCharting's card listing page was analysed using Chrome DevTools (Network → Fetch/XHR tab) to identify an internal paginated JSON endpoint returning card name, price, and daily price change data across 700+ cards in batches of 50
2. **Pagination** — the scraper increments a `cursor` parameter by 50 on each request until no products are returned, collecting all cards in one run
3. **Storage** — each card is stored in a `cards` table (one row per card), with a new row written to `price_snapshots` on every scrape — giving a time-series history per card
4. **Analysis** — the analysis package computes price movement by comparing the oldest and latest snapshot for each card, ranks by absolute percentage change, and filters by a configurable spike threshold

---

## Data Source

Prices are sourced from [PriceCharting.com](https://www.pricecharting.com/console/riftbound-origins), which aggregates Riftbound card sale data. The internal JSON endpoint used was identified through browser network traffic analysis and is not an official public API.

---

## Planned improvements

- [ ] Scheduled polling with `time.Ticker` so history builds automatically
- [ ] Discord/Slack webhook alerts on spike detection
- [ ] CSV export of price history
- [ ] Support for additional Riftbound sets (Spiritforged, Unleashed)
