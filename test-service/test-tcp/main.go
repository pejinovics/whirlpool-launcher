package main

import (
	"log"
	"net"
	"os"
	"time"
)

var startTime = time.Now()

func isAlive(elapsed time.Duration) bool {
	sec := int(elapsed.Seconds())
	switch {
	case sec < 25:
		return true
	case sec < 45:
		return false
	default:
		return true
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	addr := ":" + port
	log.Printf("TCP test service starting on %s", addr)
	log.Printf("Service will be alive for 65s, dead for 20s, then recover")

	for {
		elapsed := time.Since(startTime)
		alive := isAlive(elapsed)

		if !alive {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("[TCP] LISTENER STARTED (%v elapsed)", elapsed.Round(time.Second))
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("Failed to start listener: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			if tcpListener, ok := listener.(*net.TCPListener); ok {
				tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
			}

			conn, err := listener.Accept()

			elapsed := time.Since(startTime)
			if !isAlive(elapsed) {
				log.Printf("[TCP] LISTENER STOPPING (%v elapsed)", elapsed.Round(time.Second))
				listener.Close()
				break
			}

			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				break
			}

			conn.Write([]byte("ALIVE\n"))
			conn.Close()
			log.Printf("[TCP LIVENESS] OK (%v elapsed)", elapsed.Round(time.Second))
		}
	}
}
