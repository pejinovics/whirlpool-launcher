package helpers

import (
	"context"
	"log"

	"github.com/pejinovics/whirlpool-launcher/internal/probes/spec"
)

func CheckProbe(ctx context.Context, name, probeType string, s *spec.Specification, fail *int) bool {
	ok, _ := s.Check(ctx, s.Target)
	if ok {
		log.Printf("[%s] (%s) SUCCESS", name, probeType)
		return true
	}
	*fail++
	log.Printf("[%s] (%s) fail=%d/%d", name, probeType, *fail, s.FailureThreshold)
	return false
}
