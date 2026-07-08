package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/helpers"
	"github.com/pejinovics/whirlpool-launcher/internal/metrics"
	probesHelpers "github.com/pejinovics/whirlpool-launcher/internal/probes/helpers"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunStartup(ctx context.Context, name string, s *spec.Specification, unhealthyCh chan<- struct{}) bool {
	if s == nil {
		log.Printf("[%s] (startup) no startupProbe", name)
		return true
	}

	labels := metrics.BuildLabels(name, helpers.BuildAddress(s.Target.Host, s.Target.Port), name, "liveness")

	fail := 0
	ticker := time.NewTicker(time.Duration(s.PeriodSeconds) * time.Second)
	defer ticker.Stop()

	for {
		start := time.Now()
		ok := probesHelpers.CheckProbe(ctx, name, "startup", s, &fail)
		dur := time.Since(start)

		metrics.ObserveResult(labels, ok, dur, fail)

		if ok {
			return true
		}

		if fail >= s.FailureThreshold {
			log.Printf("[%s] (startup) FAILED", name)
			select {
			case unhealthyCh <- struct{}{}:
			default:
			}
			return false
		}

		if !helpers.WaitForNextTick(ctx, ticker) {
			return false
		}
	}
}
