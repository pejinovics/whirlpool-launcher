package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/helpers"
	"github.com/pejinovics/whirlpool-launcher/internal/metrics"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunLiveness(ctx context.Context, name string, s *spec.Specification, unhealthyCh chan<- struct{}) {
	if s == nil {
		return
	}

	if !helpers.WaitInitialDelay(ctx, s.InitialDelaySeconds) {
		return
	}

	ticker := time.NewTicker(time.Duration(s.PeriodSeconds) * time.Second)
	defer ticker.Stop()

	fail := 0
	bad := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] (liveness) stop", name)
			return
		case <-ticker.C:
			handleLivenessCheck(ctx, name, s, &fail, &bad, unhealthyCh)
		}
	}
}

func handleLivenessCheck(ctx context.Context, name string, s *spec.Specification, fail *int, bad *bool, unhealthyCh chan<- struct{}) {
	start := time.Now()
	ok, _ := s.Check(ctx, s.Target)
	dur := time.Since(start)

	labels := metrics.BuildLabels(name, helpers.BuildAddress(s.Target.Host, s.Target.Port), name, "liveness")

	if ok {
		*bad = false
		*fail = 0
		metrics.ObserveResult(labels, true, dur, *fail)
		log.Printf("[%s] (liveness) healthy", name)
		return
	}

	*fail++
	metrics.ObserveResult(labels, false, dur, *fail)
	if !*bad && *fail >= s.FailureThreshold {
		*bad = true
		log.Printf("[%s] (liveness) unhealthy", name)

		select {
		case unhealthyCh <- struct{}{}:
		default:
		}
	}
}
