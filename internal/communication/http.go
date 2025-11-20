package communication

import (
	"context"
	"fmt"
	"net/http"
)

func HTTPCheck(ctx context.Context, t Target) (bool, error) {
	timeout := getTimeout(t.Timeout)
	client := &http.Client{Timeout: timeout}

	url := buildHTTPURL(t.Host, t.Port, t.Path)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return isSuccessStatus(resp.StatusCode), nil
}

func buildHTTPURL(host string, port int, path string) string {
	return fmt.Sprintf("http://%s:%d%s", host, port, path)
}

func isSuccessStatus(code int) bool {
	return code >= 200 && code < 400
}
