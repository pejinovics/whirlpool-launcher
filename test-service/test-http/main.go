package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var startTime = time.Now()

func isReady(elapsed time.Duration) bool {
	sec := int(elapsed.Seconds())
	switch {
	case sec < 40:
		return false // još nije spreman
	case sec < 60:
		return true // postaje spreman
	case sec < 80:
		return false // opet nije spreman
	default:
		return true // posle 30s stabilno ready
	}
}

func isAlive(elapsed time.Duration) bool {
	sec := int(elapsed.Seconds())
	switch {
	case sec < 65:
		return true // živi
	case sec < 85:
		return false // crko
	default:
		return true // oporavio se
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/startup", func(w http.ResponseWriter, r *http.Request) {
		elapsed := time.Since(startTime)
		if elapsed < 30*time.Second {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf("[STARTUP] NOT READY (%v elapsed)", elapsed.Round(time.Second))
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("[STARTUP] OK (%v elapsed)", elapsed.Round(time.Second))
	})

	// LIVENESS
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		elapsed := time.Since(startTime)
		if !isAlive(elapsed) {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("[LIVENESS] FAIL (%v elapsed)", elapsed.Round(time.Second))
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("[LIVENESS] OK (%v elapsed)", elapsed.Round(time.Second))
	})

	// READINESS
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		elapsed := time.Since(startTime)
		if !isReady(elapsed) {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf("[READINESS] NOT READY (%v elapsed)", elapsed.Round(time.Second))
			return
		}
		w.WriteHeader(http.StatusOK)
		log.Printf("[READINESS] OK (%v elapsed)", elapsed.Round(time.Second))
	})

	log.Printf("Test service starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
