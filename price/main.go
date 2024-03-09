package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// PriceUpdate represents a price update event.
type PriceUpdate struct {
	Price float64 `json:"price"`
}

func main() {
	// Price update channel
	updates := make(chan float64)

	// Start a goroutine to send price updates
	go sendPriceUpdates(updates)

	// Serve SSE to clients
	http.HandleFunc("/updates", func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Listen to updates channel
		for {
			price := <-updates
			// Format the message as an SSE event
			fmt.Fprintf(w, "data: {\"price\": %f}\n\n", price)
			fmt.Println("Price updated: ", price)
			// Flush the response to send the event immediately
			w.(http.Flusher).Flush()
		}
	})

	log.Println("Server started at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// sendPriceUpdates sends price updates to the updates channel every 60 seconds.
func sendPriceUpdates(updates chan<- float64) {
	for {
		// Simulate fetching price update
		price := fetchPrice()
		updates <- price

		// Wait for 60 seconds before sending the next update
		time.Sleep(time.Second)
	}
}

// fetchPrice simulates fetching the current price.
func fetchPrice() float64 {
	// Simulate price fluctuation
	price := 100.0 + (randFloat(-5, 5))
	return price
}

// randFloat generates a random float64 number between min and max.
func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
