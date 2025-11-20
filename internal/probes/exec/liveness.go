package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunLiveness(ctx context.Context, name string, s *spec.Specification, unhealthyCh chan<- struct{}) {
	if s == nil {
		return
	}

	if s.InitialDelaySeconds > 0 {
		// log.Printf("[%s] (liveness) initialDelay=%ds", name, s.InitialDelaySeconds)
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(s.InitialDelaySeconds) * time.Second):
		}
	}

	// log.Printf("[%s] (liveness) http://%s:%d%s period=%ds FT=%d",
	// 	name, s.Target.Host, s.Target.Port, s.Target.Path,
	// 	s.PeriodSeconds, s.FailureThreshold)

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
			ok, _ := s.Check(ctx, s.Target)
			if ok {
				if bad {
					log.Printf("[%s] (liveness) RECOVERED", name)
				}
				bad = false
				fail = 0
				log.Printf("[%s] (liveness) healthy", name)
				continue
			}

			fail++
			log.Printf("[%s] (liveness) fail=%d/%d", name, fail, s.FailureThreshold)

			if !bad && fail >= s.FailureThreshold {
				bad = true
				log.Printf("[%s] (liveness) UNHEALTHY — šaljem signal orkestratoru", name)

				select {
				case unhealthyCh <- struct{}{}:
				default:
				}
			}
		}
	}
}
