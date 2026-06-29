.PHONY: build run scrape movers spikes history clean

# Build a binary you can run directly
build:
	go build -o bin/tcg ./cmd/tcg-price-tracker

# Run without building (faster for development)
run:
	go run ./cmd/tcg-price-tracker

# Shortcut commands
scrape:
	go run ./cmd/tcg-price-tracker scrape

movers:
	go run ./cmd/tcg-price-tracker movers --top 10

spikes:
	go run ./cmd/tcg-price-tracker spikes --threshold 20

history:
	go run ./cmd/tcg-price-tracker history --card "$(card)"

# Tidy dependencies
tidy:
	go mod tidy

# Delete the built binary
clean:
	rm -f bin/tcg
