package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunStartup(ctx context.Context, name string, s *spec.Specification) bool {
	if s == nil {
		log.Printf("[%s] (startup) nema startupProbe -> otključano", name)
		return true
	}

	// log.Printf("[%s] (startup) http://%s:%d%s period=%ds FT=%d",
	// 	name, s.Target.Host, s.Target.Port, s.Target.Path, s.PeriodSeconds, s.FailureThreshold)

	ticker := time.NewTicker(time.Duration(s.PeriodSeconds) * time.Second)
	defer ticker.Stop()

	fail := 0
	for {
		ok, _ := s.Check(ctx, s.Target)
		if ok {
			log.Printf("[%s] (startup) SUCCESS", name)
			return true
		}
		fail++
		log.Printf("[%s] (startup) fail=%d/%d", name, fail, s.FailureThreshold)
		if fail >= s.FailureThreshold {
			log.Printf("[%s] (startup) FAILED", name)
			return false
		}
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
		}
	}
}
