// main.go
package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	http.HandleFunc("/users/count", userCount)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalf("%+v", err)
	}
}

func userCount(w http.ResponseWriter, r *http.Request) {
	// Get the last event ID
	lastEventID, err := getEventID(*r)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	slog.Debug(fmt.Sprintf("LastEventID: %d", lastEventID))

	// Set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		lastEventID++
		// Write lines
		_, err := fmt.Fprintf(w, "id: %v\n", lastEventID) // event ID
		if err != nil {
			// TODO: Handle disconnect
			return
		}
		_, err = fmt.Fprintf(w, "retry: %d\n", 5000)       // event ID
		_, err = fmt.Fprintf(w, "event: userCount\n")      // Event name
		_, err = fmt.Fprintf(w, "data: %d\n", lastEventID) // it's not the last line
		_, err = fmt.Fprintf(w, "data: users \n\n")        // it's the last line
		// Flush
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(time.Second)
	}
}

func getEventID(r http.Request) (int64, error) {
	lastEventID := r.Header.Get("Last-Event-ID")
	if lastEventID == "" {
		return 0, nil
	}
	return strconv.ParseInt(lastEventID, 10, 64)
}
