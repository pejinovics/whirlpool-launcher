package communication

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func GRPCCheck(ctx context.Context, t Target) (bool, error) {
	to := t.Timeout
	if to <= 0 {
		to = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", t.Host, t.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := healthpb.NewHealthClient(conn)
	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{
		Service: t.Service,
	})

	if err != nil {
		return false, err
	}

	return resp.GetStatus() == healthpb.HealthCheckResponse_SERVING, nil
}
