package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Name           string `yaml:"name"`
	StartupProbe   *Probe `yaml:"startupProbe,omitempty"`
	LivenessProbe  *Probe `yaml:"livenessProbe,omitempty"`
	ReadinessProbe *Probe `yaml:"readinessProbe,omitempty"`
}

type Probe struct {
	HTTPGet             *HTTPGet   `yaml:"httpGet,omitempty"`
	TCPSocket           *TCPSocket `yaml:"tcpSocket,omitempty"`
	GRPC                *GRPC      `yaml:"grpc,omitempty"`
	InitialDelaySeconds int        `yaml:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int        `yaml:"periodSeconds"`
	FailureThreshold    int        `yaml:"failureThreshold"`
	TimeoutSeconds      int        `yaml:"timeoutSeconds,omitempty"` // default: 5s
}

type HTTPGet struct {
	Path string `yaml:"path"`
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type TCPSocket struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type GRPC struct {
	Port    int    `yaml:"port"`
	Host    string `yaml:"host,omitempty"` // default: localhost
	Service string `yaml:"service,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
