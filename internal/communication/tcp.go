package communication

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func TCPCheck(ctx context.Context, t Target) (bool, error) {
	timeout := getTimeout(t.Timeout)

	dialer := &net.Dialer{Timeout: timeout}
	addr := buildAddress(t.Host, t.Port)

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		log.Printf("[TCPCheck] ❌ Connection FAILED: %v", err)
		return false, err
	}

	conn.Close()
	log.Printf("[TCPCheck] ✅ Connection SUCCESS")

	return true, nil
}

func getTimeout(t time.Duration) time.Duration {
	if t <= 0 {
		return 5 * time.Second
	}
	return t
}

func buildAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
