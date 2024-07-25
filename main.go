package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	targetOrigin := os.Getenv("TARGET_ORIGIN")
	if targetOrigin == "" {
		log.Fatal("TARGET_ORIGIN environment variable not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		targetURL := targetOrigin + r.URL.Path
		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERROR: %s %s - %v\n", r.Method, r.URL.Path, err)
			return
		}

		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERROR: %s %s - %v\n", r.Method, r.URL.Path, err)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		duration := time.Since(startTime).Milliseconds()
		log.Printf("%s %s - %d %dms\n", r.Method, r.URL.Path, resp.StatusCode, duration)
	})

	log.Println("Starting httprelay...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
