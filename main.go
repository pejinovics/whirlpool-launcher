package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pejinovics/whirlpool-launcher/config"
	"github.com/pejinovics/whirlpool-launcher/internal/metrics"
	"github.com/pejinovics/whirlpool-launcher/internal/operations"
	"github.com/pejinovics/whirlpool-launcher/internal/runner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfgPath := "config/config.yml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	conf, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	reg := prometheus.NewRegistry()

	metrics.Register(reg)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":9090", nil)

	dockerMgr, err := operations.NewManager()
	if err != nil {
		log.Fatalf("docker manager init: %v", err)
	}
	defer dockerMgr.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	for _, c := range conf.Containers {
		c := c
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner.RunContainer(ctx, c, dockerMgr)
		}()
	}

	wg.Wait()
}
