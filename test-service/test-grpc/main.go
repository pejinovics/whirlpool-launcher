package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var startTime = time.Now()

func isReady(elapsed time.Duration) bool {
	sec := int(elapsed.Seconds())
	switch {
	case sec < 40:
		return false
	case sec < 60:
		return true
	case sec < 80:
		return false
	default:
		return true
	}
}

func isAlive(elapsed time.Duration) bool {
	sec := int(elapsed.Seconds())
	switch {
	case sec < 65:
		return true
	case sec < 85:
		return false
	default:
		return true
	}
}

type healthServer struct {
	healthpb.UnimplementedHealthServer
}

func (s *healthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	elapsed := time.Since(startTime)
	service := req.GetService()

	alive := isAlive(elapsed)
	ready := isReady(elapsed)

	var status healthpb.HealthCheckResponse_ServingStatus
	var logType string

	// Odluči status na osnovu servisa
	switch service {
	case "liveness":
		log.Printf("ovde sam sad")
		logType = "LIVENESS"
		if alive {
			status = healthpb.HealthCheckResponse_SERVING
			log.Printf("[%s] OK (%v elapsed)", logType, elapsed.Round(time.Second))
		} else {
			status = healthpb.HealthCheckResponse_NOT_SERVING
			log.Printf("[%s] FAIL (%v elapsed)", logType, elapsed.Round(time.Second))
		}

	case "readiness", "ready":
		log.Printf("ovde sam sad i ovde")
		logType = "READINESS"
		if ready {
			status = healthpb.HealthCheckResponse_SERVING
			log.Printf("[%s] OK (%v elapsed)", logType, elapsed.Round(time.Second))
		} else {
			status = healthpb.HealthCheckResponse_NOT_SERVING
			log.Printf("[%s] NOT READY (%v elapsed)", logType, elapsed.Round(time.Second))
		}

	case "startup":
		log.Printf("usaooo")
		logType = "STARTUP"
		if elapsed < 30*time.Second {
			status = healthpb.HealthCheckResponse_NOT_SERVING
			log.Printf("[%s] NOT READY (%v elapsed)", logType, elapsed.Round(time.Second))
		} else {
			status = healthpb.HealthCheckResponse_SERVING
			log.Printf("[%s] OK (%v elapsed)", logType, elapsed.Round(time.Second))
		}

	default:
		// Za nepoznate servise, koristi liveness logiku
		logType = "DEFAULT"
		if alive {
			status = healthpb.HealthCheckResponse_SERVING
			log.Printf("[%s] OK (service='%s', %v elapsed)", logType, service, elapsed.Round(time.Second))
		} else {
			status = healthpb.HealthCheckResponse_NOT_SERVING
			log.Printf("[%s] FAIL (service='%s', %v elapsed)", logType, service, elapsed.Round(time.Second))
		}
	}

	return &healthpb.HealthCheckResponse{
		Status: status,
	}, nil
}

func (s *healthServer) Watch(req *healthpb.HealthCheckRequest, stream healthpb.Health_WatchServer) error {
	// Nije implementirano za sada
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
		log.Printf("ma ovde sam")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	healthServer := &healthServer{}

	healthpb.RegisterHealthServer(grpcServer, healthServer)

	log.Printf("gRPC test service starting on :%s", port)
	log.Printf("Available health check services:")
	log.Printf("  - '' or 'liveness' or 'health' -> liveness probe")
	log.Printf("  - 'readiness' or 'ready' -> readiness probe")
	log.Printf("  - 'startup' -> startup probe")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
