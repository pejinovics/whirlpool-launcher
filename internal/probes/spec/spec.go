package spec

import (
	"context"
	"log"
	"time"

	"github.com/pejinovics/whirlpool-launcher/config"
	"github.com/pejinovics/whirlpool-launcher/internal/communication"
)

type Specification struct {
	Target              communication.Target
	Check               func(ctx context.Context, t communication.Target) (bool, error)
	InitialDelaySeconds int
	PeriodSeconds       int
	FailureThreshold    int
}

func FromConfig(p *config.Probe) *Specification {
	if p == nil {
		return nil
	}

	// timeout := time.Duration(p) * time.Second
	// if timeout == 0 {
	// 	timeout = 5 * time.Second
	// }

	timeout := 5 * time.Second // izmeniti
	target, checkFunc := CheckMethod(p, timeout)

	return &Specification{
		Target:              target,
		Check:               checkFunc,
		InitialDelaySeconds: p.InitialDelaySeconds,
		PeriodSeconds:       p.PeriodSeconds,
		FailureThreshold:    p.FailureThreshold,
	}
}

func CheckMethod(p *config.Probe, timeout time.Duration) (communication.Target, communication.CheckFunc) {
	if p.HTTPGet != nil {
		return communication.Target{
			Host:    p.HTTPGet.Host,
			Port:    p.HTTPGet.Port,
			Path:    p.HTTPGet.Path,
			Timeout: timeout,
		}, communication.HTTPCheck
	}

	if p.TCPSocket != nil {
		return communication.Target{
			Host:    p.TCPSocket.Host,
			Port:    p.TCPSocket.Port,
			Timeout: timeout,
		}, communication.TCPCheck
	}

	if p.GRPC != nil {
		log.Printf("USAOOOO ")
		log.Printf("%s", p.GRPC.Service)
		host := p.GRPC.Host
		if host == "" {
			host = "localhost"
		}
		return communication.Target{
			Host:    host,
			Port:    p.GRPC.Port,
			Service: p.GRPC.Service,
			Timeout: timeout,
		}, communication.GRPCCheck
	}

	// Ne bi trebalo da se desi zbog validacije u config.go
	panic("probe must have one check method defined")
}
