package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pejinovics/whirlpool-launcher/config"
	"github.com/pejinovics/whirlpool-launcher/internal/operations"
	"github.com/pejinovics/whirlpool-launcher/internal/runner"
)

func main() {
	// putanja do YAML-a: arg ili default
	cfgPath := "config/config.yml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	conf, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dockerMgr, err := operations.NewManager()
	if err != nil {
		log.Fatalf("docker manager init: %v", err)
	}
	defer dockerMgr.Close()

	// graceful shutdown preko signala
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	for _, c := range conf.Containers {
		c := c // capture fix
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner.RunContainer(ctx, c, dockerMgr)
		}()
	}

	log.Println("Launcher running. Press Ctrl+C to stop.")
	wg.Wait()
	log.Println("Launcher stopped.")
}
