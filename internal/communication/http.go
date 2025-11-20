package communication

import (
	"context"
	"net/http"

	"github.com/pejinovics/whirlpool-launcher/internal/helpers"
)

func HTTPCheck(ctx context.Context, t Target) (bool, error) {
	timeout := helpers.GetTimeout(t.Timeout)
	client := &http.Client{Timeout: timeout}

	url := helpers.BuildHTTPURL(t.Host, t.Port, t.Path)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return helpers.IsSuccessStatus(resp.StatusCode), nil
}
