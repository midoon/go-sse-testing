package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "error stream", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		stockChannel := make(chan string)
		go streamUpdatePrice(r.Context(), stockChannel)

		for stock := range stockChannel {
			event := fmt.Sprintf("event: %s \n"+"data: %s \n\n", "price-changed", stock)
			_, _ = fmt.Fprint(w, event)
			flusher.Flush()
		}
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func streamUpdatePrice(ctx context.Context, channel chan<- string) {
	ticker := time.NewTicker(time.Second)
	stocks := []string{"BBRI", "BBCA", "UNTR", "HEXA", "ICBP", "GOTO"}

labelBreak:
	for {
		select {
		case <-ctx.Done():
			break labelBreak
		case <-ticker.C:
			price := rand.Intn(10000)
			stockIndex := rand.Intn(len(stocks))

			channel <- fmt.Sprintf("%s-%d", stocks[stockIndex], price)
		}
	}

	ticker.Stop()
	close(channel)
}
