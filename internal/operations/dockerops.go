package operations

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Manager je tvoj "kontroler" nad Docker kontejnerima.
// Trenutno zna da radi: restart + set readiness flag (logički).
type Manager struct {
	cli *client.Client
}

// NewManager pravi Docker klijenta koristeći DOCKER_* promenljive iz okruženja,
// isto kao u zvaničnim primerima Docker Go SDK-a.
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

// // RestartContainer – bukvalno radi isto što i `docker restart <ime>`,
// // samo preko Go SDK-a.
// func (m *Manager) RestartContainer(ctx context.Context, name string) error {
// 	if m.cli == nil {
// 		return fmt.Errorf("docker client is nil")
// 	}

// 	sec := 10 // timeout u sekundama (int, ne Duration)
// 	opts := container.StopOptions{
// 		Timeout: &sec,
// 	}

// 	log.Printf("[docker] restart container=%s (timeout=%ds)", name, sec)

// 	if err := m.cli.ContainerRestart(ctx, name, opts); err != nil {
// 		return fmt.Errorf("docker restart %s: %w", name, err)
// 	}

// 	return nil
// }

func (m *Manager) RestartContainer(ctx context.Context, name string) error {
	if m.cli == nil {
		return fmt.Errorf("docker client nil")
	}

	log.Printf("[docker] stopping %s ...", name)

	// 1) Stopping (ovo koristi timeout)
	sec := 10
	stopOpts := container.StopOptions{Timeout: &sec}
	if err := m.cli.ContainerStop(ctx, name, stopOpts); err != nil {
		log.Printf("[docker] stop failed: %v", err)
		// nastavljamo – možda je već stopiran
	}

	// 2) Sačekaj malo da se oslobode resource-i (TI DEFINIŠEŠ)
	time.Sleep(5 * time.Second)

	log.Printf("[docker] starting %s ...", name)
	if err := m.cli.ContainerStart(ctx, name, container.StartOptions{}); err != nil {
		return fmt.Errorf("docker start %s: %w", name, err)
	}

	log.Printf("[docker] restart done (custom)")
	return nil
}

// SetReady – ovde NE postoji direktan Docker API koncept "ready",
// pa za sada držimo samo log + hook.
// Kasnije ovde možeš da:
//   - šalješ stanje u svoj "whirlpool" servis,
//   - upisuješ u neku bazu,
//   - menjaš labelu, itd.
func (m *Manager) SetReady(ctx context.Context, name string, ready bool) {
	log.Printf("[docker] mark container=%s ready=%v (logical readiness, not Docker native)", name, ready)
	// TODO: ovde kasnije može da ide:
	// - upis u centralni store,
	// - poziv whirlpool backenda,
	// - promene labela itd.
}
