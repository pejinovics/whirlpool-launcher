package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/helpers"
	"github.com/pejinovics/whirlpool-launcher/internal/metrics"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunReadiness(ctx context.Context, name string, s *spec.Specification) {
	if s == nil {
		return
	}

	if !helpers.WaitInitialDelay(ctx, s.InitialDelaySeconds) {
		return
	}

	ticker := time.NewTicker(time.Duration(s.PeriodSeconds) * time.Second)
	defer ticker.Stop()

	fail := 0
	notReady := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] (readiness) stop", name)
			return
		case <-ticker.C:
			handleReadinessCheck(ctx, name, s, &fail, &notReady)
		}
	}
}

func handleReadinessCheck(ctx context.Context, name string, s *spec.Specification, fail *int, notReady *bool) {
	labels := metrics.BuildLabels(name, helpers.BuildAddress(s.Target.Host, s.Target.Port), name, "liveness")
	start := time.Now()
	ok, _ := s.Check(ctx, s.Target)
	dur := time.Since(start)

	if ok {
		*notReady = false
		*fail = 0
		metrics.ObserveResult(labels, true, dur, *fail)
		log.Printf("[%s] (readiness) ready", name)
		return
	}

	*fail++
	metrics.ObserveResult(labels, false, dur, *fail)
	if !*notReady && *fail >= s.FailureThreshold {
		*notReady = true
		log.Printf("[%s] (readiness) NotReady", name)
	}
}
