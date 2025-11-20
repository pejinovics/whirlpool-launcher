package communication

import (
	"context"
	"net"

	"github.com/pejinovics/whirlpool-launcher/internal/helpers"
)

func TCPCheck(ctx context.Context, t Target) (bool, error) {
	timeout := helpers.GetTimeout(t.Timeout)

	dialer := &net.Dialer{Timeout: timeout}
	addr := helpers.BuildAddress(t.Host, t.Port)

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, err
	}

	conn.Close()
	return true, nil
}
