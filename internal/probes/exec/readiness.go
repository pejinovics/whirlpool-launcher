package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunReadiness(ctx context.Context, name string, s *spec.Specification) {
	if s == nil {
		return
	}

	if s.InitialDelaySeconds > 0 {
		// log.Printf("[%s] (readiness) initialDelay=%ds", name, s.InitialDelaySeconds)
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(s.InitialDelaySeconds) * time.Second):
		}
	}

	// log.Printf("[%s] (readiness) http://%s:%d%s period=%ds FT=%d",
	// 	name, s.Target.Host, s.Target.Port, s.Target.Path, s.PeriodSeconds, s.FailureThreshold)

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
			ok, _ := s.Check(ctx, s.Target)
			if ok {
				if notReady {
					log.Printf("[%s] (readiness) RECOVERED", name)
				}
				notReady = false
				fail = 0
				log.Printf("[%s] (readiness) ready", name)
				continue
			}
			fail++
			if !notReady && fail >= s.FailureThreshold {
				notReady = true
				log.Printf("[%s] (readiness) NotReady", name)
			}
		}
	}
}
