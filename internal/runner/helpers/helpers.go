package helpers

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/pejinovics/whirlpool-launcher/internal/operations"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/exec"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunStartupAndHandle(ctx context.Context, name string, s *spec.Specification, unhealthyCh chan<- struct{}, dm *operations.Manager) bool {
	if s == nil {
		log.Printf("[%s] (startup) no startupProbe", name)
		return true
	}

	ok := exec.RunStartup(ctx, name, s, unhealthyCh)
	if ok {
		return true
	}

	log.Printf("[%s] startup FAILED — attempting restart (if manager present)", name)
	if err := RestartContainer(ctx, dm, name); err != nil {
		log.Printf("[%s] restart attempt failed: %v", name, err)
		return false
	}

	log.Printf("[%s] restart attempt OK — caller will start new lifecycle", name)
	return false
}

func StartProbes(wg *sync.WaitGroup, ctx context.Context, name string, livenessSpec, readinessSpec *spec.Specification, unhealthyCh chan<- struct{}) {
	if livenessSpec != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exec.RunLiveness(ctx, name, livenessSpec, unhealthyCh)
		}()
	}
	if readinessSpec != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exec.RunReadiness(ctx, name, readinessSpec)
		}()
	}
}

func WaitForUnhealthyOrDone(globalCtx, genCtx context.Context, unhealthyCh <-chan struct{}) string {
	select {
	case <-globalCtx.Done():
		return "global"
	case <-genCtx.Done():
		return "global"
	case <-unhealthyCh:
		return "unhealthy"
	}
}

func RestartContainer(ctx context.Context, dm *operations.Manager, name string) error {
	if dm == nil {
		return fmt.Errorf("no Docker manager")
	}
	return dm.RestartContainer(ctx, name)
}
