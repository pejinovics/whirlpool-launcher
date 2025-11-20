package operations

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Manager struct {
	cli *client.Client
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("docker client init: %w", err)
	}
	return &Manager{cli: cli}, nil
}

func (m *Manager) Close() error {
	if m.cli != nil {
		return m.cli.Close()
	}
	return nil
}

func (m *Manager) RestartContainer(ctx context.Context, name string) error {
	if m.cli == nil {
		return fmt.Errorf("docker client is nil")
	}

	sec := 10
	opts := container.StopOptions{
		Timeout: &sec,
	}

	log.Printf("[docker] restart container=%s (timeout=%ds)", name, sec)

	if err := m.cli.ContainerRestart(ctx, name, opts); err != nil {
		return fmt.Errorf("docker restart %s: %w", name, err)
	}

	return nil
}
