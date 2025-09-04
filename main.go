package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/events", handleEvents)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	// Cấu hình SSE headers
	count := r.URL.Query().Get("count")
	if count == "" {
		count = "1"
	}
	num, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://192.168.0.103:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Channel để gửi số random
	msgChan := make(chan int)

	// Goroutine sinh dữ liệu cho client này
	go func() {
		defer close(msgChan)

		// random delay 1-5s
		time.Sleep(3*time.Second)
		msgChan <- num + 3
	}()

	// Lắng nghe context để biết khi client disconnect
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		case num, ok := <-msgChan:
			if !ok {
				return // channel đóng => dừng
			}
			// Gửi event
			fmt.Fprintf(w, "data: %d\nDone\n", num)
			flusher.Flush()
		}
	}
}
