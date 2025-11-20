package helpers

import (
	"context"
	"fmt"
	"time"
)

func BuildHTTPURL(host string, port int, path string) string {
	return fmt.Sprintf("http://%s:%d%s", host, port, path)
}

func IsSuccessStatus(code int) bool {
	return code >= 200 && code < 400
}

func GetTimeout(t time.Duration) time.Duration {
	if t <= 0 {
		return 5 * time.Second
	}
	return t
}

func BuildAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func WaitInitialDelay(ctx context.Context, delaySeconds int) bool {
	if delaySeconds > 0 {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(time.Duration(delaySeconds) * time.Second):
		}
	}
	return true
}

func WaitForNextTick(ctx context.Context, ticker *time.Ticker) bool {
	select {
	case <-ctx.Done():
		return false
	case <-ticker.C:
		return true
	}
}
