package exec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/internal/helpers"
	probesHelpers "github.com/pejinovics/whirlpool-launcher/internal/probes/helpers"
	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func RunStartup(ctx context.Context, name string, s *spec.Specification, unhealthyCh chan<- struct{}) bool {
	if s == nil {
		log.Printf("[%s] (startup) no startupProbe", name)
		return true
	}

	fail := 0
	ticker := time.NewTicker(time.Duration(s.PeriodSeconds) * time.Second)
	defer ticker.Stop()

	for {
		if probesHelpers.CheckProbe(ctx, name, "startup", s, &fail) {
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
