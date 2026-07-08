package runner

import (
	"context"
	"log"
	"sync"

	"github.com/pejinovics/whirlpool-launcher/config"
	globalhelpers "github.com/pejinovics/whirlpool-launcher/internal/helpers"
	"github.com/pejinovics/whirlpool-launcher/internal/metrics"
	"github.com/pejinovics/whirlpool-launcher/internal/operations"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
	"github.com/pejinovics/whirlpool-launcher/internal/runner/helpers"
)

func RunContainer(ctx context.Context, c config.Container, dm *operations.Manager) {
	startupSpec := spec.FromConfig(c.StartupProbe)
	livenessSpec := spec.FromConfig(c.LivenessProbe)
	readinessSpec := spec.FromConfig(c.ReadinessProbe)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] global context done", c.Name)
			return
		default:
		}

		genCtx, genCancel := context.WithCancel(ctx)
		unhealthyCh := make(chan struct{}, 1)

		startupOK := helpers.RunStartupAndHandle(genCtx, c.Name, startupSpec, unhealthyCh, dm)
		if !startupOK {
			genCancel()
			continue
		}

		var wg sync.WaitGroup
		helpers.StartProbes(&wg, genCtx, c.Name, livenessSpec, readinessSpec, unhealthyCh)

		waitStatus := helpers.WaitForUnhealthyOrDone(ctx, genCtx, unhealthyCh)
		genCancel()
		wg.Wait()

		if waitStatus == "global" {
			return
		}

		if err := helpers.RestartContainer(ctx, dm, c.Name); err != nil {
			log.Printf("[%s] Docker restart FAILED: %v", c.Name, err)
			return
		}

		labels := metrics.BuildLabels(c.Name, globalhelpers.BuildAddress(livenessSpec.Target.Host, livenessSpec.Target.Port), c.Name, "restart")
		metrics.IncRestart(labels)
	}
}
