package runner

import (
	"context"
	"log"
	"sync"

	"github.com/pejinovics/whirlpool-launcher/config"
	"github.com/pejinovics/whirlpool-launcher/internal/operations"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/exec"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunContainer(ctx context.Context, c config.Container, dm *operations.Manager) {
	log.Printf("=== Container lifecycle: %s ===", c.Name)

	// Već sada konvertujemo probe u SmallSpec, da ne radimo to u svakoj iteraciji
	startupSpec := spec.FromConfig(c.StartupProbe)
	livenessSpec := spec.FromConfig(c.LivenessProbe)
	readinessSpec := spec.FromConfig(c.ReadinessProbe)

	for {
		// Ako je globalni context prekinut (Ctrl+C / SIGTERM) — izlazimo skroz
		select {
		case <-ctx.Done():
			log.Printf("[%s] global context done -> prekid lifecycle petlje", c.Name)
			return
		default:
		}

		// 1) Novi "generation context" za ovu iteraciju (jedan cycle)
		genCtx, genCancel := context.WithCancel(ctx)

		// 2) STARTUP faza
		if ok := exec.RunStartup(genCtx, c.Name, startupSpec); !ok {
			genCancel()
			log.Printf("[%s] startup FAILED — prekidam lifecycle za ovaj container", c.Name)
			// ili umesto return možeš da radiš continue ako želiš crashloop-style
			return
		}

		// 3) LIVENESS + READINESS u gorutinama
		var wg sync.WaitGroup
		unhealthyCh := make(chan struct{}, 1)

		// Liveness – šalje signal na unhealthyCh
		if livenessSpec != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				exec.RunLiveness(genCtx, c.Name, livenessSpec, unhealthyCh)
			}()
		}

		// Readiness – može da sarađuje sa docker.Manager-om (SetReady)
		if readinessSpec != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				exec.RunReadiness(genCtx, c.Name, readinessSpec)
			}()
		}

		// 4) Čekamo ili globalni prekid, ili signal da je liveness pukao
		select {
		case <-ctx.Done():
			log.Printf("[%s] global context done u toku generation-a", c.Name)
			genCancel()
			wg.Wait()
			return

		case <-unhealthyCh:
			log.Printf("[%s] liveness javio UNHEALTHY -> gasim generation i restartujem Docker container", c.Name)
			// zaustavi ovu generaciju
			genCancel()
			wg.Wait()

			// restart Docker containera (isti naziv kao u YAML-u / Dockeru)
			if dm != nil {
				if err := dm.RestartContainer(ctx, c.Name); err != nil {
					log.Printf("[%s] Docker restart FAILED: %v", c.Name, err)
					// po želji ovde možeš da prekineš skroz:
					// return
				} else {
					log.Printf("[%s] Docker restart OK — krećemo novi lifecycle cycle", c.Name)
				}
			} else {
				log.Printf("[%s] nemam Docker manager — samo završavam ovu generaciju", c.Name)
				return
			}

			// for petlja nastavlja — ide nova iteracija:
			// novi genCtx, novi startup, itd.
		}
	}
}
