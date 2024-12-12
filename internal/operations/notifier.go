package operations

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case counts := <-sseChannel:
			jsonData, err := json.Marshal(counts)
			if err != nil {
				log.Println("Error marshalling counts to JSON:", err)
				continue
			}

			fmt.Printf("Enviando dados SSE: %+v\n", counts)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()

		case <-r.Context().Done():
			log.Println("Client disconnected")
			return
		}
	}
}
